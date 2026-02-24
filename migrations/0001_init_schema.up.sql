-- 0001_init_schema.up.sql

-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(150) NOT NULL,
    email           VARCHAR(150) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

-- Templates
CREATE TABLE templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(150) NOT NULL,
    category        VARCHAR(50) NOT NULL,
    thumbnail_url   TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_templates_category ON templates(category);

-- Template Sections
CREATE TABLE template_sections (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id     UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    default_content JSONB,
    sort_order      INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_template_sections_template_id ON template_sections(template_id);

-- Events
CREATE TABLE events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id      UUID NOT NULL REFERENCES templates(id),
    title            VARCHAR(200) NOT NULL,
    slug             VARCHAR(200) UNIQUE NOT NULL,
    event_date       TIMESTAMP NOT NULL,
    location_name    VARCHAR(255),
    location_address TEXT,
    is_published     BOOLEAN NOT NULL DEFAULT FALSE,
    view_count       INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_slug ON events(slug);

-- Event Themes
CREATE TABLE event_themes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    primary_color   VARCHAR(20),
    secondary_color VARCHAR(20),
    font_family     VARCHAR(100),
    background_url  TEXT,
    custom_css      TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_event_themes_event_id ON event_themes(event_id);

-- Event Sections
CREATE TABLE event_sections (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id            UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    template_section_id UUID NOT NULL REFERENCES template_sections(id),
    content             JSONB,
    is_visible          BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order          INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_event_sections_event_id ON event_sections(event_id);

-- Guests
CREATE TABLE guests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    phone       VARCHAR(50),
    message     TEXT,
    rsvp_status VARCHAR(20) DEFAULT 'pending',
    guest_code  VARCHAR(100),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_guests_event_id ON guests(event_id);
CREATE INDEX idx_guests_guest_code ON guests(guest_code);

-- Media
CREATE TABLE media (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID REFERENCES events(id) ON DELETE CASCADE,
    file_url    TEXT NOT NULL,
    media_type  VARCHAR(50),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_media_event_id ON media(event_id);

-- Seed: Sample templates
INSERT INTO templates (id, name, category, is_active) VALUES
    (gen_random_uuid(), 'Elegant Wedding', 'wedding', true),
    (gen_random_uuid(), 'Garden Party', 'wedding', true),
    (gen_random_uuid(), 'Birthday Blast', 'birthday', true),
    (gen_random_uuid(), 'Kids Party', 'birthday', true),
    (gen_random_uuid(), 'Community Gathering', 'community', true),
    (gen_random_uuid(), 'Aqiqah Classic', 'aqiqah', true);
