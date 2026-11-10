package handlers

import (
	"encoding/json"
	"html/template"
	"os"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// siteURL mengembalikan URL dasar situs tanpa trailing slash.
// Memakai env SITE_URL bila diatur, selain itu memakai alamat request.
func siteURL(c *fiber.Ctx) string {
	if s := strings.TrimRight(os.Getenv("SITE_URL"), "/"); s != "" {
		return s
	}
	return c.BaseURL()
}

// absURL mengubah path relatif (mis. "/public/cv.pdf") menjadi URL absolut
// yang dibutuhkan Open Graph dan JSON-LD.
func absURL(c *fiber.Ctx, u string) string {
	if u == "" || strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return siteURL(c) + u
}

// jsonLD serialisasi data schema.org menjadi JSON yang aman
// disisipkan langsung ke dalam <script type="application/ld+json">.
func jsonLD(v interface{}) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return template.JS("null")
	}
	return template.JS(b)
}

// withSEO menambahkan data canonical + Open Graph ke data render template.
func withSEO(c *fiber.Ctx, s models.SiteSetting, d fiber.Map) fiber.Map {
	d["SiteURL"] = siteURL(c)
	d["Canonical"] = siteURL(c) + c.Path()
	d["OGImage"] = absURL(c, s.ProfileImage)
	return d
}

// orString mengembalikan nilai pertama yang tidak kosong.
func orString(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// personSchema merujuk ke Person (dipakai homepage & sebagai author).
func personSchema(s models.SiteSetting, base string) map[string]interface{} {
	m := map[string]interface{}{
		"@type": "Person",
		"@id":   base + "/#person",
		"name":  orString(s.HeroTitle, s.HeroSubtitle),
		"url":   base + "/",
	}
	if s.HeroSubtitle != "" && s.HeroSubtitle != m["name"] {
		m["jobTitle"] = s.HeroSubtitle
	}
	if img := absURLAbs(base, s.ProfileImage); img != "" {
		m["image"] = img
	}
	var sameAs []string
	if s.GithubLink != "" {
		sameAs = append(sameAs, s.GithubLink)
	}
	if s.LinkedInLink != "" {
		sameAs = append(sameAs, s.LinkedInLink)
	}
	if len(sameAs) > 0 {
		m["sameAs"] = sameAs
	}
	return m
}

// absURLAbs versi personSchema yang tidak butuh request (base sudah dihitung).
func absURLAbs(base, u string) string {
	if u == "" || strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return base + u
}

// websiteSchema skema WebSite untuk homepage.
func websiteSchema(s models.SiteSetting, base string) map[string]interface{} {
	name := s.SiteTitle
	if name == "" {
		name = "Portfolio"
	}
	return map[string]interface{}{
		"@type":      "WebSite",
		"@id":        base + "/#website",
		"url":        base + "/",
		"name":       name,
		"inLanguage": "id",
		"publisher":  map[string]string{"@id": base + "/#person"},
	}
}

// HomeSchema menggabungkan Person + WebSite untuk halaman utama.
func HomeSchema(s models.SiteSetting, base string) template.JS {
	return jsonLD([]interface{}{personSchema(s, base), websiteSchema(s, base)})
}

// Sitemap menyusun sitemap.xml dinamis dari proyek & tulisan di database.
func Sitemap(c *fiber.Ctx) error {
	base := siteURL(c)

	var projects []models.Project
	database.DB.Select("slug", "created_at").Order("`order` asc, id asc").Find(&projects)

	var posts []models.BlogPost
	database.DB.Select("slug", "updated_at").Order("created_at desc").Find(&posts)

	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	writeURL := func(loc, lastmod, freq, prio string) {
		b.WriteString("  <url>\n")
		b.WriteString("    <loc>" + loc + "</loc>\n")
		if lastmod != "" {
			b.WriteString("    <lastmod>" + lastmod + "</lastmod>\n")
		}
		b.WriteString("    <changefreq>" + freq + "</changefreq>\n")
		b.WriteString("    <priority>" + prio + "</priority>\n")
		b.WriteString("  </url>\n")
	}

	today := time.Now().UTC().Format("2006-01-02")
	writeURL(base+"/", today, "daily", "1.0")
	writeURL(base+"/resume", "", "monthly", "0.8")
	writeURL(base+"/blog", "", "weekly", "0.9")
	for _, p := range projects {
		if p.Slug == "" {
			continue
		}
		writeURL(base+"/projects/"+p.Slug, p.CreatedAt.UTC().Format("2006-01-02"), "monthly", "0.7")
	}
	for _, p := range posts {
		if p.Slug == "" {
			continue
		}
		writeURL(base+"/blog/"+p.Slug, p.UpdatedAt.UTC().Format("2006-01-02"), "weekly", "0.6")
	}
	b.WriteString("</urlset>\n")

	c.Set(fiber.HeaderContentType, "application/xml; charset=utf-8")
	return c.SendString(b.String())
}

// Robots menyajikan robots.txt yang mengarahkan crawler ke sitemap.
func Robots(c *fiber.Ctx) error {
	body := "User-agent: *\n" +
		"Allow: /\n" +
		"Disallow: /admin\n" +
		"Disallow: /login\n" +
		"Disallow: /logout\n" +
		"\nSitemap: " + siteURL(c) + "/sitemap.xml\n"
	c.Set(fiber.HeaderContentType, "text/plain; charset=utf-8")
	return c.SendString(body)
}
