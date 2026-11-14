package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stream-rec-core/internal/domain"

	"github.com/go-chi/chi/v5"
)

// createItem godoc
// @Summary      Create a new catalog item
// @Description  Adds a new movie or series to the catalog
// @Tags         Items Catalog
// @Accept       json
// @Produce      json
// @Param        item  body      domain.Item  true  "Catalog item payload"
// @Success      201   {object}  domain.Item
// @Failure      500   {object}  ErrorResponse
// @Router       /items [post]
func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	if err := h.svc.CreateItem(r.Context(), &item); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusCreated, item)
}

// getItems godoc
// @Summary      List all catalog items
// @Description  Returns all movies and series stored in the catalog
// @Tags         Items Catalog
// @Produce      json
// @Success      200  {array}   domain.Item
// @Failure      500  {object}  ErrorResponse
// @Router       /items [get]
func (h *Handler) getItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetItems(r.Context())
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, items)
}

// getItemByID godoc
// @Summary      Get catalog item by ID
// @Description  Returns a single movie or series by its ObjectID
// @Tags         Items Catalog
// @Produce      json
// @Param        id   path      string  true  "Item ObjectID"
// @Success      200  {object}  domain.Item
// @Failure      404  {object}  ErrorResponse
// @Router       /items/{id} [get]
func (h *Handler) getItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.svc.GetItemByID(r.Context(), id)
	if err != nil {
		h.respond(w, http.StatusNotFound, map[string]string{"error": "item not found"})
		return
	}
	h.respond(w, http.StatusOK, item)
}

// updateItem godoc
// @Summary      Update catalog item
// @Description  Updates an existing movie or series by its ObjectID
// @Tags         Items Catalog
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Item ObjectID"
// @Param        item  body      domain.Item  true  "Updated item payload"
// @Success      200   {object}  domain.Item
// @Failure      400   {object}  ErrorResponse
// @Router       /items/{id} [put]
func (h *Handler) updateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var item domain.Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	if err := h.svc.UpdateItem(r.Context(), id, &item); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, item)
}

// deleteItem godoc
// @Summary      Delete catalog item
// @Description  Removes a movie or series from the catalog by its ObjectID
// @Tags         Items Catalog
// @Produce      json
// @Param        id   path      string  true  "Item ObjectID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  ErrorResponse
// @Router       /items/{id} [delete]
func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteItem(r.Context(), id); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// searchItems godoc
// @Summary      Search catalog items
// @Description  Finds movies or series by name substring, category, or release year
// @Tags         Items Catalog
// @Produce      json
// @Param        q         query     string  false  "Search query by item name (case-insensitive)"
// @Param        category  query     string  false  "Filter by category (e.g., Sci-Fi, Drama, Action)"
// @Param        year      query     int     false  "Filter by release year (e.g., 2014)"
// @Success      200       {array}   domain.Item
// @Failure      500       {object}  ErrorResponse
// @Router       /items/search [get]
func (h *Handler) searchItems(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	filter := domain.ItemSearchFilter{
		Query:    r.URL.Query().Get("q"),
		Category: r.URL.Query().Get("category"),
		Year:     year,
	}
	items, err := h.svc.SearchItems(r.Context(), filter)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, items)
}
