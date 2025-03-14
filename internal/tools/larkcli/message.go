package larkcli

type LarkMessage struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items []struct {
			MessageId  string `json:"message_id"`
			RootId     string `json:"root_id"`
			ParentId   string `json:"parent_id"`
			ThreadId   string `json:"thread_id"`
			MsgType    string `json:"msg_type"`
			CreateTime string `json:"create_time"`
			UpdateTime string `json:"update_time"`
			Deleted    bool   `json:"deleted"`
			Updated    bool   `json:"updated"`
			ChatId     string `json:"chat_id"`
			Sender     struct {
				Id         string `json:"id"`
				IdType     string `json:"id_type"`
				SenderType string `json:"sender_type"`
				TenantKey  string `json:"tenant_key"`
			} `json:"sender"`
			Body struct {
				Content string `json:"content"`
			} `json:"body"`
			Mentions []struct {
				Key       string `json:"key"`
				Id        string `json:"id"`
				IdType    string `json:"id_type"`
				Name      string `json:"name"`
				TenantKey string `json:"tenant_key"`
			} `json:"mentions"`
			UpperMessageId string `json:"upper_message_id"`
		} `json:"items"`
	} `json:"data"`
}
