// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/projectconfig/envmanage/env_manage.proto

package envmanage

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetEnvListReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetEnvListReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetEnvListReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetEnvListReqMultiError, or
// nil if none found.
func (m *GetEnvListReq) ValidateAll() error {
	return m.validate(true)
}

func (m *GetEnvListReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetCommon() == nil {
		err := GetEnvListReqValidationError{
			field:  "Common",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCommon()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetEnvListReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetEnvListReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommon()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetEnvListReqValidationError{
				field:  "Common",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetEnvListReqMultiError(errors)
	}

	return nil
}

// GetEnvListReqMultiError is an error wrapping multiple validation errors
// returned by GetEnvListReq.ValidateAll() if the designated constraints
// aren't met.
type GetEnvListReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetEnvListReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetEnvListReqMultiError) AllErrors() []error { return m }

// GetEnvListReqValidationError is the validation error returned by
// GetEnvListReq.Validate if the designated constraints aren't met.
type GetEnvListReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetEnvListReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetEnvListReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetEnvListReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetEnvListReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetEnvListReqValidationError) ErrorName() string { return "GetEnvListReqValidationError" }

// Error satisfies the builtin error interface
func (e GetEnvListReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetEnvListReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetEnvListReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetEnvListReqValidationError{}

// Validate checks the field values on GetEnvListRsp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetEnvListRsp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetEnvListRsp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetEnvListRspMultiError, or
// nil if none found.
func (m *GetEnvListRsp) ValidateAll() error {
	return m.validate(true)
}

func (m *GetEnvListRsp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetEnvs() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetEnvListRspValidationError{
						field:  fmt.Sprintf("Envs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetEnvListRspValidationError{
						field:  fmt.Sprintf("Envs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetEnvListRspValidationError{
					field:  fmt.Sprintf("Envs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetEnvListRspMultiError(errors)
	}

	return nil
}

// GetEnvListRspMultiError is an error wrapping multiple validation errors
// returned by GetEnvListRsp.ValidateAll() if the designated constraints
// aren't met.
type GetEnvListRspMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetEnvListRspMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetEnvListRspMultiError) AllErrors() []error { return m }

// GetEnvListRspValidationError is the validation error returned by
// GetEnvListRsp.Validate if the designated constraints aren't met.
type GetEnvListRspValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetEnvListRspValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetEnvListRspValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetEnvListRspValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetEnvListRspValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetEnvListRspValidationError) ErrorName() string { return "GetEnvListRspValidationError" }

// Error satisfies the builtin error interface
func (e GetEnvListRspValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetEnvListRsp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetEnvListRspValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetEnvListRspValidationError{}

// Validate checks the field values on Env with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Env) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Env with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EnvMultiError, or nil if none found.
func (m *Env) ValidateAll() error {
	return m.validate(true)
}

func (m *Env) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AppId

	// no validation rules for Field

	// no validation rules for Name

	// no validation rules for Remark

	// no validation rules for IsPreset

	// no validation rules for Operator

	// no validation rules for UpdateTime

	if len(errors) > 0 {
		return EnvMultiError(errors)
	}

	return nil
}

// EnvMultiError is an error wrapping multiple validation errors returned by
// Env.ValidateAll() if the designated constraints aren't met.
type EnvMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EnvMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EnvMultiError) AllErrors() []error { return m }

// EnvValidationError is the validation error returned by Env.Validate if the
// designated constraints aren't met.
type EnvValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EnvValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EnvValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EnvValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EnvValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EnvValidationError) ErrorName() string { return "EnvValidationError" }

// Error satisfies the builtin error interface
func (e EnvValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEnv.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EnvValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EnvValidationError{}

// Validate checks the field values on CreateEnvReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateEnvReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateEnvReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateEnvReqMultiError, or
// nil if none found.
func (m *CreateEnvReq) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateEnvReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetCommon() == nil {
		err := CreateEnvReqValidationError{
			field:  "Common",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCommon()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommon()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateEnvReqValidationError{
				field:  "Common",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for EnvField

	// no validation rules for EnvName

	// no validation rules for Remark

	if len(errors) > 0 {
		return CreateEnvReqMultiError(errors)
	}

	return nil
}

// CreateEnvReqMultiError is an error wrapping multiple validation errors
// returned by CreateEnvReq.ValidateAll() if the designated constraints aren't met.
type CreateEnvReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateEnvReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateEnvReqMultiError) AllErrors() []error { return m }

// CreateEnvReqValidationError is the validation error returned by
// CreateEnvReq.Validate if the designated constraints aren't met.
type CreateEnvReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateEnvReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateEnvReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateEnvReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateEnvReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateEnvReqValidationError) ErrorName() string { return "CreateEnvReqValidationError" }

// Error satisfies the builtin error interface
func (e CreateEnvReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateEnvReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateEnvReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateEnvReqValidationError{}

// Validate checks the field values on CreateEnvRsp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateEnvRsp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateEnvRsp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateEnvRspMultiError, or
// nil if none found.
func (m *CreateEnvRsp) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateEnvRsp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return CreateEnvRspMultiError(errors)
	}

	return nil
}

// CreateEnvRspMultiError is an error wrapping multiple validation errors
// returned by CreateEnvRsp.ValidateAll() if the designated constraints aren't met.
type CreateEnvRspMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateEnvRspMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateEnvRspMultiError) AllErrors() []error { return m }

// CreateEnvRspValidationError is the validation error returned by
// CreateEnvRsp.Validate if the designated constraints aren't met.
type CreateEnvRspValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateEnvRspValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateEnvRspValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateEnvRspValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateEnvRspValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateEnvRspValidationError) ErrorName() string { return "CreateEnvRspValidationError" }

// Error satisfies the builtin error interface
func (e CreateEnvRspValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateEnvRsp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateEnvRspValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateEnvRspValidationError{}

// Validate checks the field values on DeleteEnvReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DeleteEnvReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteEnvReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeleteEnvReqMultiError, or
// nil if none found.
func (m *DeleteEnvReq) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteEnvReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetCommon() == nil {
		err := DeleteEnvReqValidationError{
			field:  "Common",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCommon()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeleteEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeleteEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommon()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeleteEnvReqValidationError{
				field:  "Common",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for EnvField

	// no validation rules for EnvName

	// no validation rules for UpdateTime

	if len(errors) > 0 {
		return DeleteEnvReqMultiError(errors)
	}

	return nil
}

// DeleteEnvReqMultiError is an error wrapping multiple validation errors
// returned by DeleteEnvReq.ValidateAll() if the designated constraints aren't met.
type DeleteEnvReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteEnvReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteEnvReqMultiError) AllErrors() []error { return m }

// DeleteEnvReqValidationError is the validation error returned by
// DeleteEnvReq.Validate if the designated constraints aren't met.
type DeleteEnvReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteEnvReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteEnvReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteEnvReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteEnvReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteEnvReqValidationError) ErrorName() string { return "DeleteEnvReqValidationError" }

// Error satisfies the builtin error interface
func (e DeleteEnvReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteEnvReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteEnvReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteEnvReqValidationError{}

// Validate checks the field values on DeleteEnvRsp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DeleteEnvRsp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteEnvRsp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeleteEnvRspMultiError, or
// nil if none found.
func (m *DeleteEnvRsp) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteEnvRsp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DeleteEnvRspMultiError(errors)
	}

	return nil
}

// DeleteEnvRspMultiError is an error wrapping multiple validation errors
// returned by DeleteEnvRsp.ValidateAll() if the designated constraints aren't met.
type DeleteEnvRspMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteEnvRspMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteEnvRspMultiError) AllErrors() []error { return m }

// DeleteEnvRspValidationError is the validation error returned by
// DeleteEnvRsp.Validate if the designated constraints aren't met.
type DeleteEnvRspValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteEnvRspValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteEnvRspValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteEnvRspValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteEnvRspValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteEnvRspValidationError) ErrorName() string { return "DeleteEnvRspValidationError" }

// Error satisfies the builtin error interface
func (e DeleteEnvRspValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteEnvRsp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteEnvRspValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteEnvRspValidationError{}

// Validate checks the field values on UpdateEnvReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateEnvReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateEnvReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateEnvReqMultiError, or
// nil if none found.
func (m *UpdateEnvReq) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateEnvReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetCommon() == nil {
		err := UpdateEnvReqValidationError{
			field:  "Common",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCommon()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateEnvReqValidationError{
					field:  "Common",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommon()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateEnvReqValidationError{
				field:  "Common",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for EnvField

	// no validation rules for EnvName

	// no validation rules for NewEnvName

	// no validation rules for Remark

	// no validation rules for UpdateTime

	if len(errors) > 0 {
		return UpdateEnvReqMultiError(errors)
	}

	return nil
}

// UpdateEnvReqMultiError is an error wrapping multiple validation errors
// returned by UpdateEnvReq.ValidateAll() if the designated constraints aren't met.
type UpdateEnvReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateEnvReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateEnvReqMultiError) AllErrors() []error { return m }

// UpdateEnvReqValidationError is the validation error returned by
// UpdateEnvReq.Validate if the designated constraints aren't met.
type UpdateEnvReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateEnvReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateEnvReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateEnvReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateEnvReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateEnvReqValidationError) ErrorName() string { return "UpdateEnvReqValidationError" }

// Error satisfies the builtin error interface
func (e UpdateEnvReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateEnvReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateEnvReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateEnvReqValidationError{}

// Validate checks the field values on UpdateEnvRsp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateEnvRsp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateEnvRsp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateEnvRspMultiError, or
// nil if none found.
func (m *UpdateEnvRsp) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateEnvRsp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return UpdateEnvRspMultiError(errors)
	}

	return nil
}

// UpdateEnvRspMultiError is an error wrapping multiple validation errors
// returned by UpdateEnvRsp.ValidateAll() if the designated constraints aren't met.
type UpdateEnvRspMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateEnvRspMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateEnvRspMultiError) AllErrors() []error { return m }

// UpdateEnvRspValidationError is the validation error returned by
// UpdateEnvRsp.Validate if the designated constraints aren't met.
type UpdateEnvRspValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateEnvRspValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateEnvRspValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateEnvRspValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateEnvRspValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateEnvRspValidationError) ErrorName() string { return "UpdateEnvRspValidationError" }

// Error satisfies the builtin error interface
func (e UpdateEnvRspValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateEnvRsp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateEnvRspValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateEnvRspValidationError{}
