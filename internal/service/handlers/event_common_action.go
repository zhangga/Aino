package handlers

var (
	_ Action = (*ProcessedUniqueAction)(nil)
)

type ProcessedUniqueAction struct { //消息唯一性
}

func (ProcessedUniqueAction) Execute(data *ActionData) error {
	if data.handler.msgCache
}
