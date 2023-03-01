CREATE TABLE
  IF NOT EXISTS link_data (
    "link" text,
    "status_code" int,
    "error" text,
    "text" text,
    "images" text ARRAY,
    "audio" text ARRAY,
    "video" text ARRAY,
    "hyperlinks" text ARRAY,
    "internal_links" text ARRAY,
    "metadata" json,
    "parent_link" text,
    CONSTRAINT link_data_pkey PRIMARY KEY ("link")
  );
