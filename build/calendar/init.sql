CREATE TABLE IF NOT EXISTS Events(
    event_id SERIAL PRIMARY KEY,
    title text NOT NULL,
    date text NOT NULL,
    organizer text NOT NULL,
    public boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS Calendars(
    calendar_id SERIAL PRIMARY KEY,
    owner text NOT NULL
);

CREATE TABLE IF NOT EXISTS CalendarSharingRelationship(
    calendar_id INT NOT NULL,
    shared_calendar_id INT NOT NULL,
    PRIMARY KEY(calendar_id, shared_calendar_id),
    CONSTRAINT fk_calendar_id
        FOREIGN KEY(calendar_id)
        REFERENCES Calendars(calendar_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_shared_calendar_id
        FOREIGN KEY(shared_calendar_id)
        REFERENCES Calendars(calendar_id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CalendarEventsRelationship(
    calendar_id INT NOT NULL,
    event_id INT NOT NULL,
    PRIMARY KEY(calendar_id, event_id),
    CONSTRAINT fk_calendar_id
        FOREIGN KEY(calendar_id)
        REFERENCES Calendars(calendar_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_event_id
        FOREIGN KEY(event_id)
        REFERENCES Events(event_id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Invites(
    invite_id SERIAL PRIMARY KEY,
    status text NOT NULL,
    calendar_id INT NOT NULL,
    event_id INT NOT NULL,
    CONSTRAINT fk_calendar_id
        FOREIGN KEY(calendar_id)
        REFERENCES Calendars(calendar_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_event_id
        FOREIGN KEY(event_id)
        REFERENCES Events(event_id)
        ON DELETE CASCADE
);