// Code generated by Kitex v0.7.3. DO NOT EDIT.

package favoriteservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	favorite "wizh/kitex/kitex_gen/favorite"
)

func serviceInfo() *kitex.ServiceInfo {
	return favoriteServiceServiceInfo
}

var favoriteServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "FavoriteService"
	handlerType := (*favorite.FavoriteService)(nil)
	methods := map[string]kitex.MethodInfo{
		"FavoriteVideoAction":   kitex.NewMethodInfo(favoriteVideoActionHandler, newFavoriteVideoActionArgs, newFavoriteVideoActionResult, false),
		"FavoriteVideoList":     kitex.NewMethodInfo(favoriteVideoListHandler, newFavoriteVideoListArgs, newFavoriteVideoListResult, false),
		"FavoriteCommentAction": kitex.NewMethodInfo(favoriteCommentActionHandler, newFavoriteCommentActionArgs, newFavoriteCommentActionResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "favorite",
		"ServiceFilePath": ``,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.7.3",
		Extra:           extra,
	}
	return svcInfo
}

func favoriteVideoActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(favorite.FavoriteVideoActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(favorite.FavoriteService).FavoriteVideoAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *FavoriteVideoActionArgs:
		success, err := handler.(favorite.FavoriteService).FavoriteVideoAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FavoriteVideoActionResult)
		realResult.Success = success
	}
	return nil
}
func newFavoriteVideoActionArgs() interface{} {
	return &FavoriteVideoActionArgs{}
}

func newFavoriteVideoActionResult() interface{} {
	return &FavoriteVideoActionResult{}
}

type FavoriteVideoActionArgs struct {
	Req *favorite.FavoriteVideoActionRequest
}

func (p *FavoriteVideoActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(favorite.FavoriteVideoActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FavoriteVideoActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FavoriteVideoActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FavoriteVideoActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FavoriteVideoActionArgs) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteVideoActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FavoriteVideoActionArgs_Req_DEFAULT *favorite.FavoriteVideoActionRequest

func (p *FavoriteVideoActionArgs) GetReq() *favorite.FavoriteVideoActionRequest {
	if !p.IsSetReq() {
		return FavoriteVideoActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FavoriteVideoActionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FavoriteVideoActionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FavoriteVideoActionResult struct {
	Success *favorite.FavoriteVideoActionResponse
}

var FavoriteVideoActionResult_Success_DEFAULT *favorite.FavoriteVideoActionResponse

func (p *FavoriteVideoActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(favorite.FavoriteVideoActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FavoriteVideoActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FavoriteVideoActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FavoriteVideoActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FavoriteVideoActionResult) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteVideoActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FavoriteVideoActionResult) GetSuccess() *favorite.FavoriteVideoActionResponse {
	if !p.IsSetSuccess() {
		return FavoriteVideoActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FavoriteVideoActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*favorite.FavoriteVideoActionResponse)
}

func (p *FavoriteVideoActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FavoriteVideoActionResult) GetResult() interface{} {
	return p.Success
}

func favoriteVideoListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(favorite.FavoriteVideoListRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(favorite.FavoriteService).FavoriteVideoList(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *FavoriteVideoListArgs:
		success, err := handler.(favorite.FavoriteService).FavoriteVideoList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FavoriteVideoListResult)
		realResult.Success = success
	}
	return nil
}
func newFavoriteVideoListArgs() interface{} {
	return &FavoriteVideoListArgs{}
}

func newFavoriteVideoListResult() interface{} {
	return &FavoriteVideoListResult{}
}

type FavoriteVideoListArgs struct {
	Req *favorite.FavoriteVideoListRequest
}

func (p *FavoriteVideoListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(favorite.FavoriteVideoListRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FavoriteVideoListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FavoriteVideoListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FavoriteVideoListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FavoriteVideoListArgs) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteVideoListRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FavoriteVideoListArgs_Req_DEFAULT *favorite.FavoriteVideoListRequest

func (p *FavoriteVideoListArgs) GetReq() *favorite.FavoriteVideoListRequest {
	if !p.IsSetReq() {
		return FavoriteVideoListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FavoriteVideoListArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FavoriteVideoListArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FavoriteVideoListResult struct {
	Success *favorite.FavoriteVideoListResponse
}

var FavoriteVideoListResult_Success_DEFAULT *favorite.FavoriteVideoListResponse

func (p *FavoriteVideoListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(favorite.FavoriteVideoListResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FavoriteVideoListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FavoriteVideoListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FavoriteVideoListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FavoriteVideoListResult) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteVideoListResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FavoriteVideoListResult) GetSuccess() *favorite.FavoriteVideoListResponse {
	if !p.IsSetSuccess() {
		return FavoriteVideoListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FavoriteVideoListResult) SetSuccess(x interface{}) {
	p.Success = x.(*favorite.FavoriteVideoListResponse)
}

func (p *FavoriteVideoListResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FavoriteVideoListResult) GetResult() interface{} {
	return p.Success
}

func favoriteCommentActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(favorite.FavoriteCommentActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(favorite.FavoriteService).FavoriteCommentAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *FavoriteCommentActionArgs:
		success, err := handler.(favorite.FavoriteService).FavoriteCommentAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FavoriteCommentActionResult)
		realResult.Success = success
	}
	return nil
}
func newFavoriteCommentActionArgs() interface{} {
	return &FavoriteCommentActionArgs{}
}

func newFavoriteCommentActionResult() interface{} {
	return &FavoriteCommentActionResult{}
}

type FavoriteCommentActionArgs struct {
	Req *favorite.FavoriteCommentActionRequest
}

func (p *FavoriteCommentActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(favorite.FavoriteCommentActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FavoriteCommentActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FavoriteCommentActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FavoriteCommentActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FavoriteCommentActionArgs) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteCommentActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FavoriteCommentActionArgs_Req_DEFAULT *favorite.FavoriteCommentActionRequest

func (p *FavoriteCommentActionArgs) GetReq() *favorite.FavoriteCommentActionRequest {
	if !p.IsSetReq() {
		return FavoriteCommentActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FavoriteCommentActionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FavoriteCommentActionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FavoriteCommentActionResult struct {
	Success *favorite.FavoriteCommentActionResponse
}

var FavoriteCommentActionResult_Success_DEFAULT *favorite.FavoriteCommentActionResponse

func (p *FavoriteCommentActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(favorite.FavoriteCommentActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FavoriteCommentActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FavoriteCommentActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FavoriteCommentActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FavoriteCommentActionResult) Unmarshal(in []byte) error {
	msg := new(favorite.FavoriteCommentActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FavoriteCommentActionResult) GetSuccess() *favorite.FavoriteCommentActionResponse {
	if !p.IsSetSuccess() {
		return FavoriteCommentActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FavoriteCommentActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*favorite.FavoriteCommentActionResponse)
}

func (p *FavoriteCommentActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FavoriteCommentActionResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) FavoriteVideoAction(ctx context.Context, Req *favorite.FavoriteVideoActionRequest) (r *favorite.FavoriteVideoActionResponse, err error) {
	var _args FavoriteVideoActionArgs
	_args.Req = Req
	var _result FavoriteVideoActionResult
	if err = p.c.Call(ctx, "FavoriteVideoAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FavoriteVideoList(ctx context.Context, Req *favorite.FavoriteVideoListRequest) (r *favorite.FavoriteVideoListResponse, err error) {
	var _args FavoriteVideoListArgs
	_args.Req = Req
	var _result FavoriteVideoListResult
	if err = p.c.Call(ctx, "FavoriteVideoList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FavoriteCommentAction(ctx context.Context, Req *favorite.FavoriteCommentActionRequest) (r *favorite.FavoriteCommentActionResponse, err error) {
	var _args FavoriteCommentActionArgs
	_args.Req = Req
	var _result FavoriteCommentActionResult
	if err = p.c.Call(ctx, "FavoriteCommentAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
