package einoagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgent(ctx context.Context) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		Lambda1       = "Lambda1"
		Retriever1    = "Retriever1"
		ChatTemplate1 = "ChatTemplate1"
		Lambda2       = "Lambda2"
		Lambda3       = "Lambda3"
	)
	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(Lambda1, compose.InvokableLambda(newLambda))
	retriever1KeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever1, retriever1KeyOfRetriever)
	chatTemplate1KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate1, chatTemplate1KeyOfChatTemplate)
	lambda2KeyOfLambda, err := newLambda1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(Lambda2, lambda2KeyOfLambda)
	_ = g.AddLambdaNode(Lambda3, compose.InvokableLambda(newLambda2))
	_ = g.AddEdge(compose.START, Lambda1)
	_ = g.AddEdge(compose.START, Lambda3)
	_ = g.AddEdge(Lambda2, compose.END)
	_ = g.AddEdge(Lambda1, ChatTemplate1)
	_ = g.AddEdge(Lambda3, Retriever1)
	_ = g.AddEdge(Retriever1, ChatTemplate1)
	_ = g.AddEdge(ChatTemplate1, Lambda2)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
