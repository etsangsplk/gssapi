// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

// This file provides GSSContext methods

/*
#include <gssapi/gssapi.h>

OM_uint32
wrap_gss_init_sec_context(void *fp,
	OM_uint32 * minor_status,
	const gss_cred_id_t initiator_cred_handle,
	gss_ctx_id_t * context_handle,
	const gss_name_t target_name,
	const gss_OID mech_type,
	OM_uint32 req_flags,
	OM_uint32 time_req,
	const gss_channel_bindings_t input_chan_bindings,
	const gss_buffer_t input_token,
	gss_OID * actual_mech_type,
	gss_buffer_t output_token,
	OM_uint32 * ret_flags,
	OM_uint32 * time_rec)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		const gss_cred_id_t,
		gss_ctx_id_t *,
		const gss_name_t,
		const gss_OID,
		OM_uint32,
		OM_uint32,
		const gss_channel_bindings_t,
		const gss_buffer_t,
		gss_OID *,
		gss_buffer_t,
		OM_uint32 *,
		OM_uint32 *)
	) fp)(
		minor_status,
		initiator_cred_handle,
		context_handle,
		target_name,
		mech_type,
		req_flags,
		time_req,
		input_chan_bindings,
		input_token,
		actual_mech_type,
		output_token,
		ret_flags,
		time_rec);
}

OM_uint32
wrap_gss_accept_sec_context(void *fp,
	OM_uint32 * minor_status,
	gss_ctx_id_t * context_handle,
	const gss_cred_id_t acceptor_cred_handle,
	const gss_buffer_t input_token_buffer,
	const gss_channel_bindings_t input_chan_bindings,
	gss_name_t * src_name,
	gss_OID * mech_type,
	gss_buffer_t output_token,
	OM_uint32 * ret_flags,
	OM_uint32 * time_rec,
	gss_cred_id_t * delegated_cred_handle)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_ctx_id_t *,
		const gss_cred_id_t,
		const gss_buffer_t,
		const gss_channel_bindings_t,
		gss_name_t *,
		gss_OID *,
		gss_buffer_t,
		OM_uint32 *,
		OM_uint32 *,
		gss_cred_id_t *)
	) fp)(
		minor_status,
		context_handle,
		acceptor_cred_handle,
		input_token_buffer,
		input_chan_bindings,
		src_name,
		mech_type,
		output_token,
		ret_flags,
		time_rec,
		delegated_cred_handle);
}

OM_uint32
wrap_gss_delete_sec_context(void *fp,
	OM_uint32 * minor_status,
	gss_ctx_id_t * context_handle,
	gss_buffer_t output_token)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_ctx_id_t *,
		gss_buffer_t)
	) fp)(
		minor_status,
		context_handle,
		output_token);
}

*/
import "C"

import (
	"runtime"
	"time"
)

func (lib *Lib) NewCtxId() *CtxId {
	return &CtxId{
		Lib: lib,
	}
}

func (lib *Lib) GSS_C_NO_CONTEXT() *CtxId {
	return lib.NewCtxId()
}

// targetName is the only mandatory input?
func (lib *Lib) InitSecContext(initiatorCredHandle *CredId, ctxIn *CtxId,
	targetName *Name, mechType *OID, reqFlags uint32, timeReq time.Duration,
	inputChanBindings ChannelBindings, inputToken *Buffer) (
	ctxOut *CtxId, actualMechType *OID, outputToken *Buffer, retFlags uint32,
	timeRec time.Duration, err error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// prepare the input params
	C_initiator := C.gss_cred_id_t(nil)
	if initiatorCredHandle != nil {
		C_initiator = initiatorCredHandle.C_gss_cred_id_t
	}

	C_mechType := C.gss_OID(nil)
	if mechType != nil {
		C_mechType = mechType.C_gss_OID
	}

	C_inputToken := C.gss_buffer_t(nil)
	if inputToken != nil {
		C_inputToken = inputToken.C_gss_buffer_t
	}

	// prepare the outputs
	if ctxIn != nil {
		ctxOut = ctxIn
	} else {
		ctxOut = lib.NewCtxId()
	}

	min := C.OM_uint32(0)
	actualMechType = lib.NewOID()
	outputToken = lib.NewBuffer(true)
	flags := C.OM_uint32(0)
	timerec := C.OM_uint32(0)

	maj := C.wrap_gss_init_sec_context(lib.Fp_gss_init_sec_context,
		&min,
		C_initiator,
		&ctxOut.C_gss_ctx_id_t, // used as both in and out param
		targetName.C_gss_name_t,
		C_mechType,
		C.OM_uint32(reqFlags),
		C.OM_uint32(timeReq.Seconds()),
		C.gss_channel_bindings_t(inputChanBindings),
		C_inputToken,
		&actualMechType.C_gss_OID,
		outputToken.C_gss_buffer_t,
		&flags,
		&timerec)

	err = lib.MakeError(maj, min).GoError()
	if err != nil {
		return nil, nil, nil, 0, 0, err
	}

	return ctxOut, actualMechType, outputToken,
		uint32(flags), time.Duration(timerec) * time.Second,
		nil
}

func (lib *Lib) AcceptSecContext(
	ctxIn *CtxId, acceptorCredHandle *CredId, inputToken *Buffer,
	inputChanBindings ChannelBindings) (
	ctxOut *CtxId, srcName *Name, actualMechType *OID, outputToken *Buffer,
	retFlags uint32, timeRec time.Duration, delegatedCredHandle *CredId,
	err error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// prepare the inputs
	C_acceptorCredHandle := C.gss_cred_id_t(nil)
	if acceptorCredHandle != nil {
		C_acceptorCredHandle = acceptorCredHandle.C_gss_cred_id_t
	}

	C_inputToken := C.gss_buffer_t(nil)
	if inputToken != nil {
		C_inputToken = inputToken.C_gss_buffer_t
	}

	// prepare the outputs
	if ctxIn != nil {
		ctxOut = ctxIn
	} else {
		ctxOut = lib.GSS_C_NO_CONTEXT()
	}

	min := C.OM_uint32(0)
	srcName = lib.NewName()
	actualMechType = lib.NewOID()
	outputToken = lib.NewBuffer(true)
	flags := C.OM_uint32(0)
	timerec := C.OM_uint32(0)
	delegatedCredHandle = lib.NewCredId()

	maj := C.wrap_gss_accept_sec_context(lib.Fp_gss_accept_sec_context,
		&min,
		&ctxOut.C_gss_ctx_id_t, // used as both in and out param
		C_acceptorCredHandle,
		C_inputToken,
		C.gss_channel_bindings_t(inputChanBindings),
		&srcName.C_gss_name_t,
		&actualMechType.C_gss_OID,
		outputToken.C_gss_buffer_t,
		&flags,
		&timerec,
		&delegatedCredHandle.C_gss_cred_id_t)

	err = lib.MakeError(maj, min).GoError()
	if err != nil {
		lib.Print("AcceptSecContext: ", err)
		return nil, nil, nil, nil, 0, 0, nil, err
	}

	return ctxOut, srcName, actualMechType, outputToken, uint32(flags),
		time.Duration(timerec) * time.Second, delegatedCredHandle, nil
}

// I decided not to implement the outputToken parameter since its use is no
// longer recommended, and it would have to be Released by the caller
func (ctx *CtxId) DeleteSecContext() error {

	if ctx == nil {
		return nil
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	min := C.OM_uint32(0)
	maj := C.wrap_gss_delete_sec_context(ctx.Fp_gss_delete_sec_context,
		&min, &ctx.C_gss_ctx_id_t, nil)

	return ctx.MakeError(maj, min).GoError()
}

// TODO: gss_process_context_token
// TODO: gss_context_time
// TODO: gss_inquire_context
// TODO: gss_wrap_size_limit
// TODO: gss_export_sec_context
// TODO: gss_import_sec_context
