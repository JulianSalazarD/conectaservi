package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JulianSalazarD/conectaservi/internal/catalog"
)

// fakeCategoryRepo is an in-memory CategoryRepository for tests.
type fakeCategoryRepo struct {
	items map[uuid.UUID]*catalog.Category
}

func newFakeCategoryRepo() *fakeCategoryRepo {
	return &fakeCategoryRepo{items: map[uuid.UUID]*catalog.Category{}}
}

func (r *fakeCategoryRepo) Insert(_ context.Context, c *catalog.Category) error {
	for _, it := range r.items {
		if it.Slug == c.Slug {
			return catalog.ErrDuplicateSlug
		}
	}
	r.items[c.ID] = c
	return nil
}

func (r *fakeCategoryRepo) FindAll(_ context.Context) ([]*catalog.Category, error) {
	out := make([]*catalog.Category, 0, len(r.items))
	for _, c := range r.items {
		out = append(out, c)
	}
	return out, nil
}

func (r *fakeCategoryRepo) FindByID(_ context.Context, id uuid.UUID) (*catalog.Category, error) {
	c, ok := r.items[id]
	if !ok {
		return nil, catalog.ErrCategoryNotFound
	}
	return c, nil
}

func (r *fakeCategoryRepo) Update(_ context.Context, c *catalog.Category) error {
	if _, ok := r.items[c.ID]; !ok {
		return catalog.ErrCategoryNotFound
	}
	for id, it := range r.items {
		if id != c.ID && it.Slug == c.Slug {
			return catalog.ErrDuplicateSlug
		}
	}
	r.items[c.ID] = c
	return nil
}

func (r *fakeCategoryRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := r.items[id]; !ok {
		return catalog.ErrCategoryNotFound
	}
	delete(r.items, id)
	return nil
}

// fakeServiceRepo — minimal, only what categories tests touch.
type fakeServiceRepo struct{}

func (fakeServiceRepo) Insert(context.Context, *catalog.Service) error { return nil }
func (fakeServiceRepo) FindAll(context.Context, catalog.ServiceFilter) ([]*catalog.Service, error) {
	return nil, nil
}
func (fakeServiceRepo) FindByID(context.Context, uuid.UUID) (*catalog.Service, error) {
	return nil, catalog.ErrServiceNotFound
}
func (fakeServiceRepo) Update(context.Context, *catalog.Service) error { return nil }
func (fakeServiceRepo) Delete(context.Context, uuid.UUID) error        { return nil }

func newTestRouter(t *testing.T) (*gin.Engine, *fakeCategoryRepo) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	repo := newFakeCategoryRepo()
	mod, err := NewWithRepos(repo, fakeServiceRepo{})
	if err != nil {
		t.Fatalf("NewWithRepos: %v", err)
	}
	r := gin.New()
	mod.Mount(r)
	return r, repo
}

func TestCategoriesNewFormRendersForm(t *testing.T) {
	r, _ := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/web/categories/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `<form method="POST"`) {
		t.Errorf("expected form in body, got: %s", body)
	}
	if !strings.Contains(body, `name="nombre"`) || !strings.Contains(body, `name="slug"`) {
		t.Errorf("form fields missing")
	}
}

func TestCategoryCreateValidRedirectsToList(t *testing.T) {
	r, repo := newTestRouter(t)

	form := url.Values{}
	form.Set("nombre", "Plomería")
	form.Set("slug", "plomeria")

	req := httptest.NewRequest(http.MethodPost, "/web/categories", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d (body: %s)", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/web/categories") {
		t.Errorf("expected redirect to /web/categories, got %q", loc)
	}
	if len(repo.items) != 1 {
		t.Errorf("expected 1 category persisted, got %d", len(repo.items))
	}
}

func TestCategoryCreateInvalidSlugRerendersFormWithError(t *testing.T) {
	r, repo := newTestRouter(t)

	form := url.Values{}
	form.Set("nombre", "Plomería")
	form.Set("slug", "Plomeria!") // invalid (uppercase + symbol)

	req := httptest.NewRequest(http.MethodPost, "/web/categories", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 (re-render), got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Slug inválido") {
		t.Errorf("expected error message in HTML, got: %s", body)
	}
	if !strings.Contains(body, `value="Plomería"`) {
		t.Errorf("expected nombre value to be preserved in form")
	}
	if len(repo.items) != 0 {
		t.Errorf("expected no category persisted, got %d", len(repo.items))
	}
}

func TestCategoriesListShowsItems(t *testing.T) {
	r, repo := newTestRouter(t)
	cat, err := catalog.NewCategory("Electricidad", "electricidad", nil)
	if err != nil {
		t.Fatalf("NewCategory: %v", err)
	}
	_ = repo.Insert(context.Background(), cat)

	req := httptest.NewRequest(http.MethodGet, "/web/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Electricidad") {
		t.Errorf("expected category name in HTML")
	}
}

func TestCategoryDeleteRemovesAndRedirects(t *testing.T) {
	r, repo := newTestRouter(t)
	cat, _ := catalog.NewCategory("Limpieza", "limpieza", nil)
	_ = repo.Insert(context.Background(), cat)

	req := httptest.NewRequest(http.MethodPost, "/web/categories/"+cat.ID.String()+"/delete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if len(repo.items) != 0 {
		t.Errorf("expected 0 items after delete, got %d", len(repo.items))
	}
}

func TestHomeRenders(t *testing.T) {
	r, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "ConectaServi") {
		t.Errorf("expected brand in home")
	}
}
