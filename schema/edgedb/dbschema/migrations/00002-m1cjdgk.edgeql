CREATE MIGRATION m1cjdgkihfmeimgnlqodkcwfs6af7hxmzb3gunlo5gzts4h7ma6l5q
    ONTO m1gbwtly3ldmj6uxeea5n3fz3w2curfmw5647mwzlbvobd4j76g7cq
{
  CREATE TYPE default::Token {
      CREATE REQUIRED LINK account: default::Account;
      CREATE PROPERTY token_revoked: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY token_type: std::str {
          SET default := 'access_token';
          CREATE CONSTRAINT std::one_of('access_token', 'refresh_token');
      };
      CREATE REQUIRED PROPERTY token_value: std::str;
  };
};
