CREATE MIGRATION m1lhyel7nzow4yivkt7upm7e4adylk5aqbe5nvrs4zc3vkrymxcufa
    ONTO m16rfnep7nxq2f4k6ta7iadogiwygmhzilqiyik5fvipntqy7huptq
{
  CREATE TYPE default::Authcode {
      CREATE REQUIRED PROPERTY code: std::str;
      CREATE INDEX ON (.code);
      CREATE REQUIRED LINK account: default::Account;
      CREATE REQUIRED LINK application: default::OAuthApplication;
      CREATE REQUIRED PROPERTY consented: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY expires_at: std::datetime;
      CREATE REQUIRED PROPERTY granted_scope: array<std::str> {
          SET default := (<array<std::str>>{});
      };
      CREATE REQUIRED PROPERTY required_scope: array<std::str>;
  };
};
