// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"math/rand"
	"strings"
	"time"

	"github.com/mustang5910/deeplx-translategemma/internal/svc"
	"github.com/mustang5910/deeplx-translategemma/internal/types"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/zeromicro/go-zero/core/logx"
)

const defaultPromptTemplate = `Translate the following text into {{.TargetLang}}. Note that you should only output the translated result without any additional explanation:

{{.Text}}
`

const (
	retryBaseDelay   = 1 * time.Second
	retryMaxDelay    = 30 * time.Second
)

type TranslateLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	promptTmpl *template.Template
}

func NewTranslateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateLogic {
	promptText := svcCtx.Config.Prompt
	if promptText == "" {
		promptText = defaultPromptTemplate
	}
	return &TranslateLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		promptTmpl: template.Must(template.New("Prompt").Parse(promptText)),
	}
}

type TranslateParams struct {
	Text       string
	SourceLang string
	SourceCode string
	TargetLang string
	TargetCode string
}

func buildTranslateParams(req *types.Request) TranslateParams {
	sCode, sLang := resolveLanguageCode(req.SourceLang)
	tCode, tLang := resolveLanguageCode(req.TargetLang)

	return TranslateParams{
		Text:       req.Text,
		SourceLang: sLang,
		SourceCode: sCode,
		TargetLang: tLang,
		TargetCode: tCode,
	}
}

func resolveLanguageCode(input string) (string, string) {
	if input == "" {
		return "", ""
	}
	upper := strings.ToUpper(input)
	if info, ok := LanguageMap[upper]; ok {
		return info.Code, info.Lang
	}
	// Fallback for regions like EN-US if not explicitly key-mapped but base EN is
	if strings.Contains(upper, "-") {
		parts := strings.Split(upper, "-")
		if info, ok := LanguageMap[parts[0]]; ok {
			return info.Code, info.Lang
		}
	}
	// Final fallback
	return input, input
}

func (l *TranslateLogic) generate(translateParams TranslateParams) (string, error) {
	var promptBuffer bytes.Buffer
	if err := l.promptTmpl.Execute(&promptBuffer, translateParams); err != nil {
		return "", err
	}

	model := l.svcCtx.Config.Model
	apiKey := l.svcCtx.Config.OpenAIKey
	baseURL := l.svcCtx.Config.OpenAIBaseURL
	message := promptBuffer.String()
	maxRetries := l.svcCtx.Config.MaxRetries

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBackoff(attempt)
			l.Infof("OpenAI API retry %d/%d after %v", attempt, maxRetries, delay)
			select {
			case <-l.ctx.Done():
				return "", l.ctx.Err()
			case <-time.After(delay):
			}
		}

		chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(message),
			},
			Model: model,
		})
		if err == nil {
			return chatCompletion.Choices[0].Message.Content, nil
		}

		lastErr = err
		if !isRetryable(err) {
			return "", err
		}
	}

	return "", lastErr
}

func retryBackoff(attempt int) time.Duration {
	delay := min(retryBaseDelay*(1<<(attempt-1)), retryMaxDelay)
	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	return delay/2 + jitter
}

func isRetryable(err error) bool {
	var apiErr *openai.Error
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.StatusCode == 429:
			return true
		case apiErr.StatusCode >= 500:
			return true
		default:
			return false
		}
	}
	return true
}

func (l *TranslateLogic) Translate(req *types.Request) (resp *types.Response, err error) {
	translateParams := buildTranslateParams(req)

	// 并发控制：获取信号量槽位（超过上限则排队等待）
	if err := l.svcCtx.Semaphore.Acquire(l.ctx, 1); err != nil {
		return nil, err
	}
	defer l.svcCtx.Semaphore.Release(1)

	data, err := l.generate(translateParams)
	if err != nil {
		return nil, err
	}

	resp = &types.Response{
		Code:       200,
		Data:       data,
		SourceLang: translateParams.SourceLang,
		TargetLang: translateParams.TargetLang,
	}

	return resp, nil
}
