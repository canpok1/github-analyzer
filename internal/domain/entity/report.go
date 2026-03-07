package entity

// Report はAI分析レポートを表すエンティティ。
// Geminiの出力を構造化した各セクションを保持する。
type Report struct {
	// Overview は活動の「空気感」を1行で要約したもの。
	Overview string
	// ProcessInsights は分析指示に対する分析結果。
	ProcessInsights string
	// PotentialRisks は停滞や対立など、注意すべきリスク。
	PotentialRisks string
	// ManagersHint はマネージャーへの具体的なアクション提案。
	ManagersHint string
}
