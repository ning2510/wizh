// Code generated by Kitex v0.7.3. DO NOT EDIT.

package relationservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	relation "wizh/kitex/kitex_gen/relation"
)

func serviceInfo() *kitex.ServiceInfo {
	return relationServiceServiceInfo
}

var relationServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "RelationService"
	handlerType := (*relation.RelationService)(nil)
	methods := map[string]kitex.MethodInfo{
		"RelationAction":       kitex.NewMethodInfo(relationActionHandler, newRelationActionArgs, newRelationActionResult, false),
		"RelationFollowList":   kitex.NewMethodInfo(relationFollowListHandler, newRelationFollowListArgs, newRelationFollowListResult, false),
		"RelationFollowerList": kitex.NewMethodInfo(relationFollowerListHandler, newRelationFollowerListArgs, newRelationFollowerListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "relation",
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

func relationActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(relation.RelationActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(relation.RelationService).RelationAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *RelationActionArgs:
		success, err := handler.(relation.RelationService).RelationAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*RelationActionResult)
		realResult.Success = success
	}
	return nil
}
func newRelationActionArgs() interface{} {
	return &RelationActionArgs{}
}

func newRelationActionResult() interface{} {
	return &RelationActionResult{}
}

type RelationActionArgs struct {
	Req *relation.RelationActionRequest
}

func (p *RelationActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(relation.RelationActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *RelationActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *RelationActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *RelationActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *RelationActionArgs) Unmarshal(in []byte) error {
	msg := new(relation.RelationActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var RelationActionArgs_Req_DEFAULT *relation.RelationActionRequest

func (p *RelationActionArgs) GetReq() *relation.RelationActionRequest {
	if !p.IsSetReq() {
		return RelationActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *RelationActionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *RelationActionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type RelationActionResult struct {
	Success *relation.RelationActionResponse
}

var RelationActionResult_Success_DEFAULT *relation.RelationActionResponse

func (p *RelationActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(relation.RelationActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *RelationActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *RelationActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *RelationActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *RelationActionResult) Unmarshal(in []byte) error {
	msg := new(relation.RelationActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *RelationActionResult) GetSuccess() *relation.RelationActionResponse {
	if !p.IsSetSuccess() {
		return RelationActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *RelationActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*relation.RelationActionResponse)
}

func (p *RelationActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *RelationActionResult) GetResult() interface{} {
	return p.Success
}

func relationFollowListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(relation.RelationFollowListRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(relation.RelationService).RelationFollowList(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *RelationFollowListArgs:
		success, err := handler.(relation.RelationService).RelationFollowList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*RelationFollowListResult)
		realResult.Success = success
	}
	return nil
}
func newRelationFollowListArgs() interface{} {
	return &RelationFollowListArgs{}
}

func newRelationFollowListResult() interface{} {
	return &RelationFollowListResult{}
}

type RelationFollowListArgs struct {
	Req *relation.RelationFollowListRequest
}

func (p *RelationFollowListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(relation.RelationFollowListRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *RelationFollowListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *RelationFollowListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *RelationFollowListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *RelationFollowListArgs) Unmarshal(in []byte) error {
	msg := new(relation.RelationFollowListRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var RelationFollowListArgs_Req_DEFAULT *relation.RelationFollowListRequest

func (p *RelationFollowListArgs) GetReq() *relation.RelationFollowListRequest {
	if !p.IsSetReq() {
		return RelationFollowListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *RelationFollowListArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *RelationFollowListArgs) GetFirstArgument() interface{} {
	return p.Req
}

type RelationFollowListResult struct {
	Success *relation.RelationFollowListResponse
}

var RelationFollowListResult_Success_DEFAULT *relation.RelationFollowListResponse

func (p *RelationFollowListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(relation.RelationFollowListResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *RelationFollowListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *RelationFollowListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *RelationFollowListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *RelationFollowListResult) Unmarshal(in []byte) error {
	msg := new(relation.RelationFollowListResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *RelationFollowListResult) GetSuccess() *relation.RelationFollowListResponse {
	if !p.IsSetSuccess() {
		return RelationFollowListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *RelationFollowListResult) SetSuccess(x interface{}) {
	p.Success = x.(*relation.RelationFollowListResponse)
}

func (p *RelationFollowListResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *RelationFollowListResult) GetResult() interface{} {
	return p.Success
}

func relationFollowerListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(relation.RelationFollowerListRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(relation.RelationService).RelationFollowerList(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *RelationFollowerListArgs:
		success, err := handler.(relation.RelationService).RelationFollowerList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*RelationFollowerListResult)
		realResult.Success = success
	}
	return nil
}
func newRelationFollowerListArgs() interface{} {
	return &RelationFollowerListArgs{}
}

func newRelationFollowerListResult() interface{} {
	return &RelationFollowerListResult{}
}

type RelationFollowerListArgs struct {
	Req *relation.RelationFollowerListRequest
}

func (p *RelationFollowerListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(relation.RelationFollowerListRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *RelationFollowerListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *RelationFollowerListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *RelationFollowerListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *RelationFollowerListArgs) Unmarshal(in []byte) error {
	msg := new(relation.RelationFollowerListRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var RelationFollowerListArgs_Req_DEFAULT *relation.RelationFollowerListRequest

func (p *RelationFollowerListArgs) GetReq() *relation.RelationFollowerListRequest {
	if !p.IsSetReq() {
		return RelationFollowerListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *RelationFollowerListArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *RelationFollowerListArgs) GetFirstArgument() interface{} {
	return p.Req
}

type RelationFollowerListResult struct {
	Success *relation.RelationFollowerListResponse
}

var RelationFollowerListResult_Success_DEFAULT *relation.RelationFollowerListResponse

func (p *RelationFollowerListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(relation.RelationFollowerListResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *RelationFollowerListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *RelationFollowerListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *RelationFollowerListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *RelationFollowerListResult) Unmarshal(in []byte) error {
	msg := new(relation.RelationFollowerListResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *RelationFollowerListResult) GetSuccess() *relation.RelationFollowerListResponse {
	if !p.IsSetSuccess() {
		return RelationFollowerListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *RelationFollowerListResult) SetSuccess(x interface{}) {
	p.Success = x.(*relation.RelationFollowerListResponse)
}

func (p *RelationFollowerListResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *RelationFollowerListResult) GetResult() interface{} {
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

func (p *kClient) RelationAction(ctx context.Context, Req *relation.RelationActionRequest) (r *relation.RelationActionResponse, err error) {
	var _args RelationActionArgs
	_args.Req = Req
	var _result RelationActionResult
	if err = p.c.Call(ctx, "RelationAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RelationFollowList(ctx context.Context, Req *relation.RelationFollowListRequest) (r *relation.RelationFollowListResponse, err error) {
	var _args RelationFollowListArgs
	_args.Req = Req
	var _result RelationFollowListResult
	if err = p.c.Call(ctx, "RelationFollowList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RelationFollowerList(ctx context.Context, Req *relation.RelationFollowerListRequest) (r *relation.RelationFollowerListResponse, err error) {
	var _args RelationFollowerListArgs
	_args.Req = Req
	var _result RelationFollowerListResult
	if err = p.c.Call(ctx, "RelationFollowerList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
