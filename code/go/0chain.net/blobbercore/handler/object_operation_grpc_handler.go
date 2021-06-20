package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/allocation"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/blobbergrpc"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/convert"
	"github.com/0chain/blobber/code/go/0chain.net/core/common"
	"mime/multipart"
	"net/http"
)

func (b *blobberGRPCService) UpdateObjectAttributes(ctx context.Context, req *blobbergrpc.UpdateObjectAttributesRequest) (*blobbergrpc.UpdateObjectAttributesResponse, error) {
	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		return nil, err
	}
	httpRequestWithMetaData(r, GetGRPCMetaDataFromCtx(ctx), req.Allocation)
	r.Form = map[string][]string{
		"path":          {req.Path},
		"path_hash":     {req.PathHash},
		"connection_id": {req.ConnectionId},
		"attributes":    {req.Attributes},
	}

	resp, err := UpdateAttributesHandler(ctx, r)
	if err != nil {
		return nil, err
	}

	return convert.UpdateObjectAttributesResponseCreator(resp), nil
}

func (b *blobberGRPCService) CopyObject(ctx context.Context, req *blobbergrpc.CopyObjectRequest) (*blobbergrpc.CopyObjectResponse, error) {
	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		return nil, err
	}
	httpRequestWithMetaData(r, GetGRPCMetaDataFromCtx(ctx), req.Allocation)
	r.Form = map[string][]string{
		"path":          {req.Path},
		"path_hash":     {req.PathHash},
		"connection_id": {req.ConnectionId},
		"dest":          {req.Dest},
	}

	resp, err := CopyHandler(ctx, r)
	if err != nil {
		return nil, err
	}

	return convert.CopyObjectResponseCreator(resp), nil
}

func (b *blobberGRPCService) RenameObject(ctx context.Context, req *blobbergrpc.RenameObjectRequest) (*blobbergrpc.RenameObjectResponse, error) {
	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		return nil, err
	}
	httpRequestWithMetaData(r, GetGRPCMetaDataFromCtx(ctx), req.Allocation)
	r.Form = map[string][]string{
		"path":          {req.Path},
		"path_hash":     {req.PathHash},
		"connection_id": {req.ConnectionId},
		"new_name":      {req.NewName},
	}

	resp, err := RenameHandler(ctx, r)
	if err != nil {
		return nil, err
	}

	return convert.RenameObjectResponseCreator(resp), nil
}

func (b *blobberGRPCService) DownloadFile(ctx context.Context, req *blobbergrpc.DownloadFileRequest) (*blobbergrpc.DownloadFileResponse, error) {

	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		return nil, err
	}

	httpRequestWithMetaData(r, GetGRPCMetaDataFromCtx(ctx), req.Allocation)
	r.Form = map[string][]string{
		"path":        {req.Path},
		"path_hash":   {req.PathHash},
		"rx_pay":      {req.RxPay},
		"block_num":   {req.BlockNum},
		"num_blocks":  {req.NumBlocks},
		"read_marker": {req.ReadMarker},
		"auth_token":  {req.AuthToken},
		"content":     {req.AuthToken},
	}

	resp, err := DownloadHandler(ctx, r)
	if err != nil {
		return nil, err
	}

	return convert.DownloadFileResponseCreator(resp), nil
}

func (b *blobberGRPCService) WriteFile(ctx context.Context, req *blobbergrpc.UploadFileRequest) (*blobbergrpc.UploadFileResponse, error) {

	var formData allocation.UpdateFileChange
	var uploadMetaString string
	switch req.Method {
	case `POST`:
		uploadMetaString = req.UploadMeta
	case `PUT`:
		uploadMetaString = req.UpdateMeta
	}
	err := json.Unmarshal([]byte(uploadMetaString), &formData)
	if err != nil {
		return nil, common.NewError("invalid_parameters",
			"Invalid parameters. Error parsing the meta data for upload."+err.Error())
	}

	r, err := http.NewRequest(req.Method, "", nil)
	if err != nil {
		return nil, err
	}

	if req.Method != `DELETE` {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(`uploadFile`, formData.Filename)
		if err != nil {
			return nil, err
		}
		_, err = part.Write(req.UploadFile)
		if err != nil {
			return nil, err
		}

		thumbPart, err := writer.CreateFormFile(`uploadThumbnailFile`, formData.ThumbnailFilename)
		if err != nil {
			return nil, err
		}
		_, err = thumbPart.Write(req.UploadThumbnailFile)
		if err != nil {
			return nil, err
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		r, err = http.NewRequest(req.Method, "", body)
		if err != nil {
			return nil, err
		}
	}

	httpRequestWithMetaData(r, GetGRPCMetaDataFromCtx(ctx), req.Allocation)
	r.Form = map[string][]string{
		"path":          {req.Path},
		"connection_id": {req.ConnectionId},
		"uploadMeta":    {req.UploadMeta},
		"updateMeta":    {req.UpdateMeta},
	}

	resp, err := UploadHandler(ctx, r)
	if err != nil {
		return nil, err
	}

	return convert.UploadFileResponseCreator(resp), nil
}
