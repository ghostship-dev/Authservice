CREATE MIGRATION m1izmclnysxwtzdusxdnf7o4gsvqx3ps6p66jnckvawg2hfslf6fja
    ONTO m1ebs7a66ao6qmaeo2yo6p3w6bwnujdn3hliyy2eehlp4f6sh24drq
{
  CREATE TYPE default::OAuthApplication {
      CREATE REQUIRED PROPERTY client_id: std::str;
      CREATE INDEX ON (.client_id);
      CREATE REQUIRED LINK client_owner: default::Account;
      CREATE PROPERTY client_description: std::str;
      CREATE PROPERTY client_homepage_url: std::str;
      CREATE PROPERTY client_logo_url: std::str;
      CREATE REQUIRED PROPERTY client_name: std::str;
      CREATE PROPERTY client_privacy_url: std::str;
      CREATE PROPERTY client_rate_limits: std::json;
      CREATE REQUIRED PROPERTY client_registration_date: std::datetime;
      CREATE REQUIRED PROPERTY client_secret: std::str;
      CREATE REQUIRED PROPERTY client_status: std::str {
          SET default := 'enabled';
          CREATE CONSTRAINT std::one_of('active', 'disabled', 'suspended');
      };
      CREATE PROPERTY client_tos_url: std::str;
      CREATE REQUIRED PROPERTY client_type: std::str;
      CREATE REQUIRED PROPERTY grant_types: array<std::str> {
          SET default := (['authorization_code']);
      };
      CREATE REQUIRED PROPERTY redirect_uris: array<std::str> {
          SET default := (<array<std::str>>{});
      };
      CREATE REQUIRED PROPERTY scope: array<std::str> {
          SET default := (<array<std::str>>{});
      };
  };
};
