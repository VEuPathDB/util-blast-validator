package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Foxcapades/go-midl/v2/pkg/midl"
	"github.com/francoispqt/gojay"
	"github.com/gorilla/mux"
	"github.com/veupathdb/lib-go-blast/v2/pkg/blast"
	"github.com/veupathdb/lib-go-blast/v2/pkg/blast/bval"
)

const (
	Path  = "/validate/{tool}"
	param = "tool"
)

func RegisterEndpoint(router *mux.Router) {
	router.Handle(Path, midl.JSONAdapter(midl.MiddlewareFunc(handle))).
		Methods(http.MethodPost)
}

type response struct {
	Status  uint16               `json:"status"`
	Message string               `json:"message,omitempty"`
	Payload bval.ValidationError `json:"payload,omitempty"`
}

func handle(req midl.Request) midl.Response {
	tool := blast.Tool(mux.Vars(req.RawRequest())[param])

	log.Println("Validating config for tool ", tool)

	var tmp gojay.UnmarshalerJSONObject

	switch tool {
	case blast.ToolBlastN:
		tmp = &blast.BlastN{}
	case blast.ToolBlastP:
		tmp = &blast.BlastP{}
	case blast.ToolBlastX:
		tmp = &blast.BlastX{}
	case blast.ToolDeltaBlast:
		tmp = &blast.DeltaBlast{}
	case blast.ToolPSIBlast:
		tmp = &blast.PSIBlast{}
	case blast.ToolRPSBlast:
		tmp = &blast.RPSBlast{}
	case blast.ToolRPSTBlastN:
		tmp = &blast.RPSTBlastN{}
	case blast.ToolTBlastN:
		tmp = &blast.TBlastN{}
	case blast.ToolTBlastX:
		tmp = &blast.TBlastX{}
	case blast.ToolBlastFormatter:
		tmp = &blast.BlastFormatter{}
	default:
		log.Println("Unrecognized tool, returning 404.")
		return midl.MakeResponse(404, response{
			Status:  404,
			Message: fmt.Sprintf(`"%s" is not a recognized blast+ tool`, tool),
		})
	}

	if err := gojay.UnmarshalJSONObject(req.Body(), tmp); err != nil {
		log.Println("Failed to parse JSON input.")
		return midl.MakeResponse(400, response{
			Status:  400,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	log.Println("Validation completed.")
	return midl.MakeResponse(200, response{
		Status:  200,
		Payload: tmp.(bval.Validator).Validate(),
	})
}
