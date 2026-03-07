package gemini

// geminiRequest はGemini APIへのリクエストボディ。
type geminiRequest struct {
	Contents []content `json:"contents"`
}

// content はリクエスト/レスポンスのコンテンツ。
type content struct {
	Parts []part `json:"parts"`
}

// part はコンテンツの一部（テキスト）。
type part struct {
	Text string `json:"text"`
}

// geminiResponse はGemini APIの成功レスポンス。
type geminiResponse struct {
	Candidates []candidate `json:"candidates"`
}

// candidate はレスポンスの候補。
type candidate struct {
	Content content `json:"content"`
}

// geminiErrorResponse はGemini APIのエラーレスポンス。
type geminiErrorResponse struct {
	Error geminiError `json:"error"`
}

// geminiError はエラー詳細。
type geminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
