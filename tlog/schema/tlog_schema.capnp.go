// Code generated by capnpc-go. DO NOT EDIT.

package schema

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type TlogResponse struct{ capnp.Struct }

// TlogResponse_TypeID is the unique identifier for the type TlogResponse.
const TlogResponse_TypeID = 0x98d11ae1c78a24d9

func NewTlogResponse(s *capnp.Segment) (TlogResponse, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return TlogResponse{st}, err
}

func NewRootTlogResponse(s *capnp.Segment) (TlogResponse, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return TlogResponse{st}, err
}

func ReadRootTlogResponse(msg *capnp.Message) (TlogResponse, error) {
	root, err := msg.RootPtr()
	return TlogResponse{root.Struct()}, err
}

func (s TlogResponse) String() string {
	str, _ := text.Marshal(0x98d11ae1c78a24d9, s.Struct)
	return str
}

func (s TlogResponse) Status() int8 {
	return int8(s.Struct.Uint8(0))
}

func (s TlogResponse) SetStatus(v int8) {
	s.Struct.SetUint8(0, uint8(v))
}

func (s TlogResponse) Sequences() (capnp.UInt64List, error) {
	p, err := s.Struct.Ptr(0)
	return capnp.UInt64List{List: p.List()}, err
}

func (s TlogResponse) HasSequences() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TlogResponse) SetSequences(v capnp.UInt64List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewSequences sets the sequences field to a newly
// allocated capnp.UInt64List, preferring placement in s's segment.
func (s TlogResponse) NewSequences(n int32) (capnp.UInt64List, error) {
	l, err := capnp.NewUInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.UInt64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// TlogResponse_List is a list of TlogResponse.
type TlogResponse_List struct{ capnp.List }

// NewTlogResponse creates a new list of TlogResponse.
func NewTlogResponse_List(s *capnp.Segment, sz int32) (TlogResponse_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return TlogResponse_List{l}, err
}

func (s TlogResponse_List) At(i int) TlogResponse { return TlogResponse{s.List.Struct(i)} }

func (s TlogResponse_List) Set(i int, v TlogResponse) error { return s.List.SetStruct(i, v.Struct) }

// TlogResponse_Promise is a wrapper for a TlogResponse promised by a client call.
type TlogResponse_Promise struct{ *capnp.Pipeline }

func (p TlogResponse_Promise) Struct() (TlogResponse, error) {
	s, err := p.Pipeline.Struct()
	return TlogResponse{s}, err
}

type TlogBlock struct{ capnp.Struct }

// TlogBlock_TypeID is the unique identifier for the type TlogBlock.
const TlogBlock_TypeID = 0x8cf178de3c82d431

func NewTlogBlock(s *capnp.Segment) (TlogBlock, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 3})
	return TlogBlock{st}, err
}

func NewRootTlogBlock(s *capnp.Segment) (TlogBlock, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 3})
	return TlogBlock{st}, err
}

func ReadRootTlogBlock(msg *capnp.Message) (TlogBlock, error) {
	root, err := msg.RootPtr()
	return TlogBlock{root.Struct()}, err
}

func (s TlogBlock) String() string {
	str, _ := text.Marshal(0x8cf178de3c82d431, s.Struct)
	return str
}

func (s TlogBlock) VdiskID() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s TlogBlock) HasVdiskID() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TlogBlock) VdiskIDBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s TlogBlock) SetVdiskID(v string) error {
	return s.Struct.SetText(0, v)
}

func (s TlogBlock) Sequence() uint64 {
	return s.Struct.Uint64(0)
}

func (s TlogBlock) SetSequence(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s TlogBlock) Lba() uint64 {
	return s.Struct.Uint64(8)
}

func (s TlogBlock) SetLba(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s TlogBlock) Size() uint64 {
	return s.Struct.Uint64(16)
}

func (s TlogBlock) SetSize(v uint64) {
	s.Struct.SetUint64(16, v)
}

func (s TlogBlock) Hash() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return []byte(p.Data()), err
}

func (s TlogBlock) HasHash() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s TlogBlock) SetHash(v []byte) error {
	return s.Struct.SetData(1, v)
}

func (s TlogBlock) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return []byte(p.Data()), err
}

func (s TlogBlock) HasData() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s TlogBlock) SetData(v []byte) error {
	return s.Struct.SetData(2, v)
}

func (s TlogBlock) Timestamp() uint64 {
	return s.Struct.Uint64(24)
}

func (s TlogBlock) SetTimestamp(v uint64) {
	s.Struct.SetUint64(24, v)
}

func (s TlogBlock) Operation() uint8 {
	return s.Struct.Uint8(32)
}

func (s TlogBlock) SetOperation(v uint8) {
	s.Struct.SetUint8(32, v)
}

// TlogBlock_List is a list of TlogBlock.
type TlogBlock_List struct{ capnp.List }

// NewTlogBlock creates a new list of TlogBlock.
func NewTlogBlock_List(s *capnp.Segment, sz int32) (TlogBlock_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 40, PointerCount: 3}, sz)
	return TlogBlock_List{l}, err
}

func (s TlogBlock_List) At(i int) TlogBlock { return TlogBlock{s.List.Struct(i)} }

func (s TlogBlock_List) Set(i int, v TlogBlock) error { return s.List.SetStruct(i, v.Struct) }

// TlogBlock_Promise is a wrapper for a TlogBlock promised by a client call.
type TlogBlock_Promise struct{ *capnp.Pipeline }

func (p TlogBlock_Promise) Struct() (TlogBlock, error) {
	s, err := p.Pipeline.Struct()
	return TlogBlock{s}, err
}

type TlogAggregation struct{ capnp.Struct }

// TlogAggregation_TypeID is the unique identifier for the type TlogAggregation.
const TlogAggregation_TypeID = 0xe46ab5b4b619e094

func NewTlogAggregation(s *capnp.Segment) (TlogAggregation, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 4})
	return TlogAggregation{st}, err
}

func NewRootTlogAggregation(s *capnp.Segment) (TlogAggregation, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 4})
	return TlogAggregation{st}, err
}

func ReadRootTlogAggregation(msg *capnp.Message) (TlogAggregation, error) {
	root, err := msg.RootPtr()
	return TlogAggregation{root.Struct()}, err
}

func (s TlogAggregation) String() string {
	str, _ := text.Marshal(0xe46ab5b4b619e094, s.Struct)
	return str
}

func (s TlogAggregation) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s TlogAggregation) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s TlogAggregation) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s TlogAggregation) Size() uint64 {
	return s.Struct.Uint64(0)
}

func (s TlogAggregation) SetSize(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s TlogAggregation) Timestamp() uint64 {
	return s.Struct.Uint64(8)
}

func (s TlogAggregation) SetTimestamp(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s TlogAggregation) VdiskID() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s TlogAggregation) HasVdiskID() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) VdiskIDBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s TlogAggregation) SetVdiskID(v string) error {
	return s.Struct.SetText(1, v)
}

func (s TlogAggregation) Blocks() (TlogBlock_List, error) {
	p, err := s.Struct.Ptr(2)
	return TlogBlock_List{List: p.List()}, err
}

func (s TlogAggregation) HasBlocks() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) SetBlocks(v TlogBlock_List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewBlocks sets the blocks field to a newly
// allocated TlogBlock_List, preferring placement in s's segment.
func (s TlogAggregation) NewBlocks(n int32) (TlogBlock_List, error) {
	l, err := NewTlogBlock_List(s.Struct.Segment(), n)
	if err != nil {
		return TlogBlock_List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

func (s TlogAggregation) Prev() ([]byte, error) {
	p, err := s.Struct.Ptr(3)
	return []byte(p.Data()), err
}

func (s TlogAggregation) HasPrev() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) SetPrev(v []byte) error {
	return s.Struct.SetData(3, v)
}

// TlogAggregation_List is a list of TlogAggregation.
type TlogAggregation_List struct{ capnp.List }

// NewTlogAggregation creates a new list of TlogAggregation.
func NewTlogAggregation_List(s *capnp.Segment, sz int32) (TlogAggregation_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 4}, sz)
	return TlogAggregation_List{l}, err
}

func (s TlogAggregation_List) At(i int) TlogAggregation { return TlogAggregation{s.List.Struct(i)} }

func (s TlogAggregation_List) Set(i int, v TlogAggregation) error {
	return s.List.SetStruct(i, v.Struct)
}

// TlogAggregation_Promise is a wrapper for a TlogAggregation promised by a client call.
type TlogAggregation_Promise struct{ *capnp.Pipeline }

func (p TlogAggregation_Promise) Struct() (TlogAggregation, error) {
	s, err := p.Pipeline.Struct()
	return TlogAggregation{s}, err
}

const schema_f4533cbae6e08506 = "x\xda|\x93AH\x14m\x1c\xc6\x9f\xe7\xff\xdfY?" +
	"A\xddo\xd89|\xca\x17Ix1\xb2\xb4n\"X" +
	"\xd2\xa1<\xf9Z\xd01\xc6uX7ww6g\xb4" +
	"\xe8\x12\x04\x1e\x8a\x8e\x1e\xec\x10\x18x((\xe8P\x1d" +
	"\xa2K\x10\xd8\xa1\x8b``P\xb1\x82E\x81A\x81A" +
	"\x9d&^\x97]\xb7\x90n\xef\xfb\xf0\xcc\xfc\x9f\xff\xf3" +
	"\x9b\xe9\x9f\xe4q\x19pB\x01\xcc\xffN:\x19X\xbb" +
	"6\xf4\xfe\xf2\xb7\x9b0\x9dt\x92\xf4|\xf5\xe3\xd3\xa1" +
	"3\xdbp\xb4\x058\xb6\x8f]\xcc\xf6\xd1\x1e{\xb9B" +
	"0y\xd3sce\xa3ku\xd1\xda\xd9d\xdf\xf1\xfc" +
	"\x94\xa3\xcc\xb6\xda'\xb3\x8e^\x02\x93\x85j\xe7\x93G" +
	"\x8f/lZ\xb74\xb9S\xd6=\xab\xa3\xcc^\xdf\x99" +
	"3\xaf\xe7\x88\xbe$\xcaM\x05%\xffH\xac\xc50\x7f" +
	"\xbev9\x9c\xf3+\xe5\xca\xe0\xd9b\x98\x1f)\x86\x9a" +
	"\x9b\x1e#M\xb7\xa6\x80\x14\x01wu\x040\xaf\x94f" +
	"]Hz\xb4\xda\xebQ\xc0\xac)MU\xe8\x0a=\x0a" +
	"\xe0\xbe;\x00\x98u\xa5\xd9\x14\xba*\x1e\x15p7\x0e" +
	"\x02\xe6\xad\xd2|\x12\xba)zL\x01\xee\x07+V\x95" +
	"fK\xe8:\xe2\xd1\x01\xdc\xcfV\xdcT\x9a\xafB7" +
	"\xad\x1e\xd3\x80\xfbe\x1c0[J\xf3C\xe8\xb6t{" +
	"\xb6\x03\xf7\xbb\x15\xb7\x95\xe3\x14^\x9d\x9b,D\xd3\xa7" +
	"O\xb2\x0d\xc260\x89\x82\x8b\xb3A9\x17\x00`+" +
	"\x84\xad`Kq\xc2\xaf\x9f3Q\xe1J\xd0\xb8L\xf9" +
	"\xd1\x14\xdb!l\x073\x93~\xec\xd7/I\\(\x05" +
	"Q\xec\x97\xc0J\xdd\x9d\x84\x95`\xc6\x8f\x0b!Xf" +
	"\x1a\xc2\xb4\x1d\xf7\xd7>\xc7\x83h\x7f%,G\x81\xad" +
	"\xf4\x9fF\xa5\xbd\x83\x80\xe9Q\x9a\xfe\xddJ\xfb\xecR" +
	"\x87\x94\xe6\x94p8\x8a\xfdx6\xa2@(\xcd;1" +
	"b\x078\xa6\xdc\xc9\xd4\xd14?\xb5\xe7\xfc\x13\xf9\xfc" +
	"L\x90\xb7\x99\xcb\x80\xcd\xf0_#\xc3-[\xf7\x82\xd2" +
	",\xedf\xb8m\xb5E\xa5Yn\xc2z\xc7\x06[R" +
	"\x9a\xfb\x16+kX\xef\xd9\x8fbYi\x1eZ\xacR" +
	"\xc3\xfa\xc0\xaeuWi^X\xacZ\xc3\xfa\xdc\xbe\xf3" +
	"\x99\xd2\xbc\x14f\xca~)\xa8\x93\xfa\x8d\xc4^}\xff" +
	"\x89vx\xa2\x18\xe6\xa6\x1b\x0d\xfc\xbb\xfb_\x81V\xcc" +
	"Tf\x82\xb9:\xc0_\x01\x00\x00\xff\xff\xf0v\xc7\x08"

func init() {
	schemas.Register(schema_f4533cbae6e08506,
		0x8cf178de3c82d431,
		0x98d11ae1c78a24d9,
		0xe46ab5b4b619e094)
}
