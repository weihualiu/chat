package server

// 聊天连接信息及通道
//type ChatMsg struct {
//	id    string      //序列唯一标识
//	msgCh chan []byte //对应的channel
//}

type ChatContainer map[string]chan []byte

var chatContainer ChatContainer
