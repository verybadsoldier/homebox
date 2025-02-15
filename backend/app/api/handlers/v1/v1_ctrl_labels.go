package v1

import (
	"net/http"

	"github.com/hay-kot/homebox/backend/internal/core/services"
	"github.com/hay-kot/homebox/backend/internal/data/ent"
	"github.com/hay-kot/homebox/backend/internal/data/repo"
	"github.com/hay-kot/homebox/backend/internal/sys/validate"
	"github.com/hay-kot/homebox/backend/pkgs/server"
	"github.com/rs/zerolog/log"
)

// HandleLabelsGetAll godoc
//
//	@Summary  Get All Labels
//	@Tags     Labels
//	@Produce  json
//	@Success  200 {object} server.Results{items=[]repo.LabelOut}
//	@Router   /v1/labels [GET]
//	@Security Bearer
func (ctrl *V1Controller) HandleLabelsGetAll() server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		user := services.UseUserCtx(r.Context())
		labels, err := ctrl.repo.Labels.GetAll(r.Context(), user.GroupID)
		if err != nil {
			log.Err(err).Msg("error getting labels")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}
		return server.Respond(w, http.StatusOK, server.Results{Items: labels})
	}
}

// HandleLabelsCreate godoc
//
//	@Summary  Create Label
//	@Tags     Labels
//	@Produce  json
//	@Param    payload body     repo.LabelCreate true "Label Data"
//	@Success  200     {object} repo.LabelSummary
//	@Router   /v1/labels [POST]
//	@Security Bearer
func (ctrl *V1Controller) HandleLabelsCreate() server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		createData := repo.LabelCreate{}
		if err := server.Decode(r, &createData); err != nil {
			log.Err(err).Msg("error decoding label create data")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		user := services.UseUserCtx(r.Context())
		label, err := ctrl.repo.Labels.Create(r.Context(), user.GroupID, createData)
		if err != nil {
			log.Err(err).Msg("error creating label")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.Respond(w, http.StatusCreated, label)
	}
}

// HandleLabelDelete godocs
//
//	@Summary  Delete Label
//	@Tags     Labels
//	@Produce  json
//	@Param    id path string true "Label ID"
//	@Success  204
//	@Router   /v1/labels/{id} [DELETE]
//	@Security Bearer
func (ctrl *V1Controller) HandleLabelDelete() server.HandlerFunc {
	return ctrl.handleLabelsGeneral()
}

// HandleLabelGet godocs
//
//	@Summary  Get Label
//	@Tags     Labels
//	@Produce  json
//	@Param    id  path     string true "Label ID"
//	@Success  200 {object} repo.LabelOut
//	@Router   /v1/labels/{id} [GET]
//	@Security Bearer
func (ctrl *V1Controller) HandleLabelGet() server.HandlerFunc {
	return ctrl.handleLabelsGeneral()
}

// HandleLabelUpdate godocs
//
//	@Summary  Update Label
//	@Tags     Labels
//	@Produce  json
//	@Param    id  path     string true "Label ID"
//	@Success  200 {object} repo.LabelOut
//	@Router   /v1/labels/{id} [PUT]
//	@Security Bearer
func (ctrl *V1Controller) HandleLabelUpdate() server.HandlerFunc {
	return ctrl.handleLabelsGeneral()
}

func (ctrl *V1Controller) handleLabelsGeneral() server.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := services.NewContext(r.Context())
		ID, err := ctrl.routeID(r)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			labels, err := ctrl.repo.Labels.GetOneByGroup(r.Context(), ctx.GID, ID)
			if err != nil {
				if ent.IsNotFound(err) {
					log.Err(err).
						Str("id", ID.String()).
						Msg("label not found")
					return validate.NewRequestError(err, http.StatusNotFound)
				}
				log.Err(err).Msg("error getting label")
				return validate.NewRequestError(err, http.StatusInternalServerError)
			}
			return server.Respond(w, http.StatusOK, labels)

		case http.MethodDelete:
			err = ctrl.repo.Labels.DeleteByGroup(ctx, ctx.GID, ID)
			if err != nil {
				log.Err(err).Msg("error deleting label")
				return validate.NewRequestError(err, http.StatusInternalServerError)
			}
			return server.Respond(w, http.StatusNoContent, nil)

		case http.MethodPut:
			body := repo.LabelUpdate{}
			if err := server.Decode(r, &body); err != nil {
				return validate.NewRequestError(err, http.StatusInternalServerError)
			}

			body.ID = ID
			result, err := ctrl.repo.Labels.UpdateByGroup(ctx, ctx.GID, body)
			if err != nil {
				return validate.NewRequestError(err, http.StatusInternalServerError)
			}
			return server.Respond(w, http.StatusOK, result)
		}

		return nil
	}
}
