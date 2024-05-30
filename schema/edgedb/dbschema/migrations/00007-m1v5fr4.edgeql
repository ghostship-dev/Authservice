CREATE MIGRATION m1v5fr4o34ukw26znuyxv5fg5hlevr7nx4c534zazmcwkcqwue42ga
    ONTO m1ksoa72vlh22hosxarsrds6m2sibj5upghn45xq5f4ipcvroc7fmq
{
  ALTER TYPE default::Token {
      DROP INDEX ON (.token_value);
  };
  ALTER TYPE default::Token {
      ALTER PROPERTY token_value {
          RENAME TO value;
      };
  };
  ALTER TYPE default::Token {
      CREATE INDEX ON (.value);
  };
  ALTER TYPE default::Token {
      CREATE REQUIRED PROPERTY scope: array<std::str> {
          SET REQUIRED USING (<array<std::str>>{});
      };
  };
  ALTER TYPE default::Token {
      ALTER PROPERTY token_revoked {
          RENAME TO revoked;
      };
  };
  ALTER TYPE default::Token {
      ALTER PROPERTY token_type {
          RENAME TO variant;
      };
  };
};
