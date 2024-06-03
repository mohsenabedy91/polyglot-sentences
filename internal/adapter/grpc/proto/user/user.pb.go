// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: internal/adapter/grpc/proto/user/user.proto

package user

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for GetByUUID.
type GetByUUIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique UUID of the user.
	UserUUID string `protobuf:"bytes,1,opt,name=userUUID,proto3" json:"userUUID,omitempty"`
}

func (x *GetByUUIDRequest) Reset() {
	*x = GetByUUIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByUUIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByUUIDRequest) ProtoMessage() {}

func (x *GetByUUIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByUUIDRequest.ProtoReflect.Descriptor instead.
func (*GetByUUIDRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{0}
}

func (x *GetByUUIDRequest) GetUserUUID() string {
	if x != nil {
		return x.UserUUID
	}
	return ""
}

// Request message for GetByEmail.
type GetByEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The email of the user.
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *GetByEmailRequest) Reset() {
	*x = GetByEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByEmailRequest) ProtoMessage() {}

func (x *GetByEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByEmailRequest.ProtoReflect.Descriptor instead.
func (*GetByEmailRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *GetByEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// Request message for IsEmailUnique.
type IsEmailUniqueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The email to check for uniqueness.
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *IsEmailUniqueRequest) Reset() {
	*x = IsEmailUniqueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsEmailUniqueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsEmailUniqueRequest) ProtoMessage() {}

func (x *IsEmailUniqueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsEmailUniqueRequest.ProtoReflect.Descriptor instead.
func (*IsEmailUniqueRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{2}
}

func (x *IsEmailUniqueRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// Request message for UpdateGoogleID.
type UpdateGoogleIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user ID for whom to update google id.
	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	// The google identifier user have a request register/login.
	GoogleId string `protobuf:"bytes,2,opt,name=googleId,proto3" json:"googleId,omitempty"`
}

func (x *UpdateGoogleIDRequest) Reset() {
	*x = UpdateGoogleIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateGoogleIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGoogleIDRequest) ProtoMessage() {}

func (x *UpdateGoogleIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGoogleIDRequest.ProtoReflect.Descriptor instead.
func (*UpdateGoogleIDRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateGoogleIDRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UpdateGoogleIDRequest) GetGoogleId() string {
	if x != nil {
		return x.GoogleId
	}
	return ""
}

// Request message for UpdateLastLoginTime.
type UpdateLastLoginTimeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user ID for whom to update last login time.
	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *UpdateLastLoginTimeRequest) Reset() {
	*x = UpdateLastLoginTimeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLastLoginTimeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLastLoginTimeRequest) ProtoMessage() {}

func (x *UpdateLastLoginTimeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLastLoginTimeRequest.ProtoReflect.Descriptor instead.
func (*UpdateLastLoginTimeRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateLastLoginTimeRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// Request message for UpdatePassword.
type UpdatePasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user ID for whom to update password.
	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	// The hashed password of the user, there is a secure value.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *UpdatePasswordRequest) Reset() {
	*x = UpdatePasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePasswordRequest) ProtoMessage() {}

func (x *UpdatePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePasswordRequest.ProtoReflect.Descriptor instead.
func (*UpdatePasswordRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{5}
}

func (x *UpdatePasswordRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UpdatePasswordRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Request message for Create.
type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The universally unique identifier of the user.
	UUID string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	// The first name of the user.
	FirstName string `protobuf:"bytes,2,opt,name=firstName,proto3" json:"firstName,omitempty"`
	// The last name of the user.
	LastName string `protobuf:"bytes,3,opt,name=lastName,proto3" json:"lastName,omitempty"`
	// The email of the user.
	Email string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	// The hashed password of the user, there is a secure value.
	Password string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	// The avatar url of the user.
	Avatar string `protobuf:"bytes,6,opt,name=avatar,proto3" json:"avatar,omitempty"`
	// The google id of the user.
	GoogleId string `protobuf:"bytes,7,opt,name=googleId,proto3" json:"googleId,omitempty"`
	// The status of the user.
	Status string `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{6}
}

func (x *CreateRequest) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *CreateRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *CreateRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *CreateRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *CreateRequest) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *CreateRequest) GetGoogleId() string {
	if x != nil {
		return x.GoogleId
	}
	return ""
}

func (x *CreateRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// Request message for VerifiedEmail.
type VerifiedEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The email to mark as verified.
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *VerifiedEmailRequest) Reset() {
	*x = VerifiedEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifiedEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifiedEmailRequest) ProtoMessage() {}

func (x *VerifiedEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifiedEmailRequest.ProtoReflect.Descriptor instead.
func (*VerifiedEmailRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{7}
}

func (x *VerifiedEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// Request message for MarkWelcomeMessageSent.
type UpdateWelcomeMessageToSentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user ID for whom to update the welcome message sent flag.
	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *UpdateWelcomeMessageToSentRequest) Reset() {
	*x = UpdateWelcomeMessageToSentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateWelcomeMessageToSentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateWelcomeMessageToSentRequest) ProtoMessage() {}

func (x *UpdateWelcomeMessageToSentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateWelcomeMessageToSentRequest.ProtoReflect.Descriptor instead.
func (*UpdateWelcomeMessageToSentRequest) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateWelcomeMessageToSentRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// Response message containing user details.
type UserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID of the user.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The unique UUID of the user.
	UUID string `protobuf:"bytes,2,opt,name=UUID,proto3" json:"UUID,omitempty"`
	// The first name of the user.
	FirstName string `protobuf:"bytes,3,opt,name=firstName,proto3" json:"firstName,omitempty"`
	// The last name of the user.
	LastName string `protobuf:"bytes,4,opt,name=lastName,proto3" json:"lastName,omitempty"`
	// The email of the user.
	Email string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	// The status of the user.
	Status string `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	// The hashed password of the user, there is a secure value.
	Password string `protobuf:"bytes,7,opt,name=password,proto3" json:"password,omitempty"`
	// Whether the welcome message has been sent.
	WelcomeMessageSent bool `protobuf:"varint,8,opt,name=welcomeMessageSent,proto3" json:"welcomeMessageSent,omitempty"`
	// The google Id of user has a authentication request.
	GoogleId string `protobuf:"bytes,9,opt,name=googleId,proto3" json:"googleId,omitempty"`
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_adapter_grpc_proto_user_user_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP(), []int{9}
}

func (x *UserResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserResponse) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *UserResponse) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UserResponse) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UserResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UserResponse) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *UserResponse) GetWelcomeMessageSent() bool {
	if x != nil {
		return x.WelcomeMessageSent
	}
	return false
}

func (x *UserResponse) GetGoogleId() string {
	if x != nil {
		return x.GoogleId
	}
	return ""
}

var File_internal_adapter_grpc_proto_user_user_proto protoreflect.FileDescriptor

var file_internal_adapter_grpc_proto_user_user_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x64, 0x61, 0x70, 0x74,
	0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2e, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x79, 0x55, 0x55, 0x49, 0x44, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x55, 0x55, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x55, 0x55, 0x49, 0x44,
	0x22, 0x29, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x2c, 0x0a, 0x14, 0x49,
	0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x4b, 0x0a, 0x15, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x34, 0x0a, 0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4b, 0x0a, 0x15,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xdb, 0x01, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x55,
	0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12,
	0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2c, 0x0a, 0x14, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x3b, 0x0a, 0x21, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57,
	0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x53,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x82, 0x02, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x2e, 0x0a, 0x12, 0x77,
	0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x6e,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x64, 0x32, 0xf8, 0x04, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x79,
	0x55, 0x55, 0x49, 0x44, 0x12, 0x16, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x42,
	0x79, 0x55, 0x55, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x39, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x17,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0d, 0x49,
	0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x12, 0x1a, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x49, 0x73, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x55, 0x6e, 0x69, 0x71, 0x75,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x31, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x12, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x59, 0x0a, 0x16, 0x4d, 0x61, 0x72, 0x6b,
	0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65,
	0x6e, 0x74, 0x12, 0x27, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x53, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x45, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x1b, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4f, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x20, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c,
	0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x45, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1b, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x42, 0x4e, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6d, 0x6f, 0x68, 0x73, 0x65, 0x6e, 0x61, 0x62, 0x65, 0x64, 0x79, 0x39, 0x31, 0x2f, 0x70,
	0x6f, 0x6c, 0x79, 0x67, 0x6c, 0x6f, 0x74, 0x2d, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x64, 0x61, 0x70, 0x74,
	0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_adapter_grpc_proto_user_user_proto_rawDescOnce sync.Once
	file_internal_adapter_grpc_proto_user_user_proto_rawDescData = file_internal_adapter_grpc_proto_user_user_proto_rawDesc
)

func file_internal_adapter_grpc_proto_user_user_proto_rawDescGZIP() []byte {
	file_internal_adapter_grpc_proto_user_user_proto_rawDescOnce.Do(func() {
		file_internal_adapter_grpc_proto_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_adapter_grpc_proto_user_user_proto_rawDescData)
	})
	return file_internal_adapter_grpc_proto_user_user_proto_rawDescData
}

var file_internal_adapter_grpc_proto_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_internal_adapter_grpc_proto_user_user_proto_goTypes = []interface{}{
	(*GetByUUIDRequest)(nil),                  // 0: user.GetByUUIDRequest
	(*GetByEmailRequest)(nil),                 // 1: user.GetByEmailRequest
	(*IsEmailUniqueRequest)(nil),              // 2: user.IsEmailUniqueRequest
	(*UpdateGoogleIDRequest)(nil),             // 3: user.UpdateGoogleIDRequest
	(*UpdateLastLoginTimeRequest)(nil),        // 4: user.UpdateLastLoginTimeRequest
	(*UpdatePasswordRequest)(nil),             // 5: user.UpdatePasswordRequest
	(*CreateRequest)(nil),                     // 6: user.CreateRequest
	(*VerifiedEmailRequest)(nil),              // 7: user.VerifiedEmailRequest
	(*UpdateWelcomeMessageToSentRequest)(nil), // 8: user.UpdateWelcomeMessageToSentRequest
	(*UserResponse)(nil),                      // 9: user.UserResponse
	(*empty.Empty)(nil),                       // 10: google.protobuf.Empty
}
var file_internal_adapter_grpc_proto_user_user_proto_depIdxs = []int32{
	0,  // 0: user.UserService.GetByUUID:input_type -> user.GetByUUIDRequest
	1,  // 1: user.UserService.GetByEmail:input_type -> user.GetByEmailRequest
	2,  // 2: user.UserService.IsEmailUnique:input_type -> user.IsEmailUniqueRequest
	6,  // 3: user.UserService.Create:input_type -> user.CreateRequest
	7,  // 4: user.UserService.VerifiedEmail:input_type -> user.VerifiedEmailRequest
	8,  // 5: user.UserService.MarkWelcomeMessageSent:input_type -> user.UpdateWelcomeMessageToSentRequest
	3,  // 6: user.UserService.UpdateGoogleID:input_type -> user.UpdateGoogleIDRequest
	4,  // 7: user.UserService.UpdateLastLoginTime:input_type -> user.UpdateLastLoginTimeRequest
	5,  // 8: user.UserService.UpdatePassword:input_type -> user.UpdatePasswordRequest
	9,  // 9: user.UserService.GetByUUID:output_type -> user.UserResponse
	9,  // 10: user.UserService.GetByEmail:output_type -> user.UserResponse
	10, // 11: user.UserService.IsEmailUnique:output_type -> google.protobuf.Empty
	9,  // 12: user.UserService.Create:output_type -> user.UserResponse
	10, // 13: user.UserService.VerifiedEmail:output_type -> google.protobuf.Empty
	10, // 14: user.UserService.MarkWelcomeMessageSent:output_type -> google.protobuf.Empty
	10, // 15: user.UserService.UpdateGoogleID:output_type -> google.protobuf.Empty
	10, // 16: user.UserService.UpdateLastLoginTime:output_type -> google.protobuf.Empty
	10, // 17: user.UserService.UpdatePassword:output_type -> google.protobuf.Empty
	9,  // [9:18] is the sub-list for method output_type
	0,  // [0:9] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_internal_adapter_grpc_proto_user_user_proto_init() }
func file_internal_adapter_grpc_proto_user_user_proto_init() {
	if File_internal_adapter_grpc_proto_user_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByUUIDRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByEmailRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsEmailUniqueRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateGoogleIDRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLastLoginTimeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePasswordRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifiedEmailRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateWelcomeMessageToSentRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_adapter_grpc_proto_user_user_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_adapter_grpc_proto_user_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_adapter_grpc_proto_user_user_proto_goTypes,
		DependencyIndexes: file_internal_adapter_grpc_proto_user_user_proto_depIdxs,
		MessageInfos:      file_internal_adapter_grpc_proto_user_user_proto_msgTypes,
	}.Build()
	File_internal_adapter_grpc_proto_user_user_proto = out.File
	file_internal_adapter_grpc_proto_user_user_proto_rawDesc = nil
	file_internal_adapter_grpc_proto_user_user_proto_goTypes = nil
	file_internal_adapter_grpc_proto_user_user_proto_depIdxs = nil
}
