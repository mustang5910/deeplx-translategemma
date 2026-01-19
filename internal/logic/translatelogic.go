// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"bytes"
	"context"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/mustang5910/deeplx-translategemma/internal/svc"
	"github.com/mustang5910/deeplx-translategemma/internal/types"
	"github.com/ollama/ollama/api"

	"github.com/zeromicro/go-zero/core/logx"
)

var promptTmpl *template.Template

func init() {
	const templateText = `You are a professional {{.SourceLang}} ({{.SourceCode}}) to {{.TargetLang}} ({{.TargetCode}}) translator. Your goal is to accurately convey the meaning and nuances of the original {{.SourceLang}} text while adhering to {{.TargetLang}} grammar, vocabulary, and cultural sensitivities.
Produce only the {{.TargetLang}} translation, without any additional explanations or commentary.
CRITICAL: All HTML-style tags (including custom or numbered tags like <i1>, <i2>, etc.) must be preserved EXACTLY as they are in the source text. Do NOT change tag names, do NOT replace them with standard tags, and do NOT remove them. Only translate the text content between the tags.
CRITICAL: All HTML-style tags (e.g., <i1>, <i2>) must be preserved EXACTLY. Do NOT change tag names, do NOT remove tags, and do NOT add any new tags.
Example:
Input: <i1>Hello World</l1>
Correct Output: <i1>你好世界</l1>
Incorrect Output: <你好世界> or <i1>你好世界</i1>
CRITICAL: Do NOT add any extra text, labels, or context. The output must contain ONLY the translated content within the exact same tag structure as the source.
CRITICAL: If a word or phrase has multiple valid translations, you must choose the single most contextually appropriate one. Do NOT provide synonyms, do NOT list alternatives, and do NOT use separators like "/" or "or" to show options.
Please translate the following {{.SourceLang}} text into {{.TargetLang}}:




{{.Text}}`
	promptTmpl = template.Must(template.New("Prompt").Parse(templateText))
}

type TranslateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTranslateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateLogic {
	return &TranslateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
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

func (l *TranslateLogic) generate_by_ollama(translateParams TranslateParams) (string, error) {
	var promptBuffer bytes.Buffer
	if err := promptTmpl.Execute(&promptBuffer, translateParams); err != nil {
		return "", err
	}
	ollamaUrl, err := url.Parse(l.svcCtx.Config.OllamaUrl)
	if err != nil {
		return "", err
	}
	client := api.NewClient(ollamaUrl, http.DefaultClient)
	request := &api.GenerateRequest{
		Model:  l.svcCtx.Config.Model,
		Prompt: promptBuffer.String(),

		// set streaming to false
		Stream: new(bool),
	}

	value := ""

	err = client.Generate(l.ctx, request, func(resp api.GenerateResponse) error {
		value = resp.Response
		return nil
	})
	if err != nil {
		return "", err
	}

	return value, nil
}

func (l *TranslateLogic) Translate(req *types.Request) (resp *types.Response, err error) {
	translateParams := buildTranslateParams(req)

	data, err := l.generate_by_ollama(translateParams)
	if err != nil {
		return nil, err
	}

	// Fix common LLM output issues for tags in Chinese/HTML context
	data = html.UnescapeString(data)
	data = strings.ReplaceAll(data, "＜", "<")
	data = strings.ReplaceAll(data, "＞", ">")
	data = strings.ReplaceAll(data, "〈", "<")
	data = strings.ReplaceAll(data, "〉", ">")
	data = strings.ReplaceAll(data, "《", "<")
	data = strings.ReplaceAll(data, "》", ">")

	resp = &types.Response{
		Code:       200,
		Data:       data,
		SourceLang: translateParams.SourceLang,
		TargetLang: translateParams.TargetLang,
	}

	return resp, nil
}
