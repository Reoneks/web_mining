CREATE TABLE
  IF NOT EXISTS sites (
    "domain_id" text,
    "base_url" text,
    "punycode" text,
    "dns_sec" boolean,
    "name_servers" text ARRAY,
    "status" text ARRAY,
    "whois_server" text,
    "created_date" text,
    "updated_date" text,
    "expiration_date" text,
    -----------------------------------
    "exclude" text ARRAY,
    "headers" json,
    CONSTRAINT sites_pkey PRIMARY KEY ("base_url")
  );
