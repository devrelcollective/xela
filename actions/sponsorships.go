package actions

import (
	"github.com/PagerDuty/xela/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Sponsorship)
// DB Table: Plural (sponsorships)
// Resource: Plural (Sponsorships)
// Path: Plural (/sponsorships)
// View Template Folder: Plural (/templates/sponsorships/)

// SponsorshipsResource is the resource for the Sponsorship model
type SponsorshipsResource struct {
	buffalo.Resource
}

// List gets all Sponsorships. This function is mapped to the path
// GET /sponsorships
func (v SponsorshipsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	sponsorships := &models.Sponsorships{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Sponsorships from the DB
	if err := q.All(sponsorships); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, sponsorships))
}

// Show gets the data for one Sponsorship. This function is mapped to
// the path GET /sponsorships/{sponsorship_id}
func (v SponsorshipsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sponsorship
	sponsorship := &models.Sponsorship{}

	// To find the Sponsorship the parameter sponsorship_id is used.
	if err := tx.Find(sponsorship, c.Param("sponsorship_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, sponsorship))
}

// New renders the form for creating a new Sponsorship.
// This function is mapped to the path GET /sponsorships/new
func (v SponsorshipsResource) New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	events := models.Events{}
	if err := tx.Order("title").All(&events); err != nil {
		return errors.WithStack(err)
	}
	sponsorship := &models.Sponsorship{}
	c.Set("sponsorship", sponsorship)
	c.Set("events", events)
	return c.Render(422, r.Auto(c, sponsorship))
}

// Create adds a Sponsorship to the DB. This function is mapped to the
// path POST /sponsorships
func (v SponsorshipsResource) Create(c buffalo.Context) error {

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	user := c.Value("current_user").(*models.User)

	// Allocate an empty Sponsorship
	sponsorship := &models.Sponsorship{
		UserID:    user.ID,
		UpdatedBy: user.ID,
	}

	// Bind sponsorship to the html form elements
	if err := c.Bind(sponsorship); err != nil {
		return errors.WithStack(err)
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(sponsorship)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, sponsorship))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Sponsorship was created successfully")

	// and redirect to the sponsorships index page
	return c.Render(201, r.Auto(c, sponsorship))
}

// Edit renders a edit form for a Sponsorship. This function is
// mapped to the path GET /sponsorships/{sponsorship_id}/edit
func (v SponsorshipsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	events := models.Events{}
	if err := tx.Order("title").All(&events); err != nil {
		return errors.WithStack(err)
	}

	sponsorship := &models.Sponsorship{}
	if err := tx.Find(sponsorship, c.Param("sponsorship_id")); err != nil {
		return c.Error(404, err)
	}
	c.Set("sponsorship", sponsorship)
	c.Set("events", events)
	return c.Render(200, r.Auto(c, sponsorship))
}

// Update changes a Sponsorship in the DB. This function is mapped to
// the path PUT /sponsorships/{sponsorship_id}
func (v SponsorshipsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	user := c.Value("current_user").(*models.User)

	// Allocate an empty Sponsorship
	sponsorship := &models.Sponsorship{
		UpdatedBy: user.ID,
	}

	if err := tx.Find(sponsorship, c.Param("sponsorship_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Sponsorship to the html form elements
	if err := c.Bind(sponsorship); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(sponsorship)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, sponsorship))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Sponsorship was updated successfully")

	// and redirect to the sponsorships index page
	return c.Render(200, r.Auto(c, sponsorship))
}

// Destroy deletes a Sponsorship from the DB. This function is mapped
// to the path DELETE /sponsorships/{sponsorship_id}
func (v SponsorshipsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sponsorship
	sponsorship := &models.Sponsorship{}

	// To find the Sponsorship the parameter sponsorship_id is used.
	if err := tx.Find(sponsorship, c.Param("sponsorship_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(sponsorship); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Sponsorship was destroyed successfully")

	// Redirect to the sponsorships index page
	return c.Render(200, r.Auto(c, sponsorship))
}