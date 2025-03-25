package dto

// CreateQuestionDTO 创建问题的数据传输对象
type CreateQuestionDTO struct {
	Title          string   // 标题
	Content        []byte   // HyperText 对象
	SimpleQuestion string   // 简单文本内容
	Type           string   // 题目类型：单选、多选、填空等
	Difficulty     string   // 难度等级：使用预定义的难度常量，如 DifficultyCefrA1, DifficultyHsk1 等
	Options        []byte   // 选项列表
	OptionTuples   []byte   // 选项双元组列表
	Answers        []string // 答案列表
	Category       string   // 题目分类
	Labels         []string // 标签列表
	Explanation    string   // 解析
	Attachments    []string // 附件列表
	TimeLimit      uint32   // 时间限制，单位秒
}

// NewCreateQuestionDTO 创建新的 CreateQuestionDTO 实例
func NewCreateQuestionDTO(
	title string,
	content []byte,
	simpleQuestion string,
	questionType string,
	difficulty string, // 使用预定义的难度常量，如 DifficultyCefrA1, DifficultyHsk1 等
	options []byte,
	optionTuples []byte,
	answers []string,
	category string,
	labels []string,
	explanation string,
	attachments []string,
	timeLimit uint32,
) *CreateQuestionDTO {
	return &CreateQuestionDTO{
		Title:          title,
		Content:        content,
		SimpleQuestion: simpleQuestion,
		Type:           questionType,
		Difficulty:     difficulty,
		Options:        options,
		OptionTuples:   optionTuples,
		Answers:        answers,
		Category:       category,
		Labels:         labels,
		Explanation:    explanation,
		Attachments:    attachments,
		TimeLimit:      timeLimit,
	}
}

// UpdateQuestionDTO 更新问题的数据传输对象
type UpdateQuestionDTO struct {
	ID             string   // 问题ID
	Title          string   // 标题
	Content        []byte   // HyperText 对象
	SimpleQuestion string   // 简单文本内容
	Type           string   // 题目类型：单选、多选、填空等
	Difficulty     string   // 难度等级：使用预定义的难度常量，如 DifficultyCefrA1, DifficultyHsk1 等
	Options        []byte   // 选项列表
	OptionTuples   []byte   // 选项双元组列表
	Answers        []string // 答案列表
	Category       string   // 题目分类
	Labels         []string // 标签列表
	Explanation    string   // 解析
	Attachments    []string // 附件列表
	TimeLimit      uint32   // 时间限制，单位秒
}

// NewUpdateQuestionDTO 创建新的 UpdateQuestionDTO 实例
func NewUpdateQuestionDTO(
	id string,
	title string,
	content []byte,
	simpleQuestion string,
	questionType string,
	difficulty string,
	options []byte,
	optionTuples []byte,
	answers []string,
	category string,
	labels []string,
	explanation string,
	attachments []string,
	timeLimit uint32,
) *UpdateQuestionDTO {
	return &UpdateQuestionDTO{
		ID:             id,
		Title:          title,
		Content:        content,
		SimpleQuestion: simpleQuestion,
		Type:           questionType,
		Difficulty:     difficulty,
		Options:        options,
		OptionTuples:   optionTuples,
		Answers:        answers,
		Category:       category,
		Labels:         labels,
		Explanation:    explanation,
		Attachments:    attachments,
		TimeLimit:      timeLimit,
	}
}
