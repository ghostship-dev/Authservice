CREATE MIGRATION m1gbwtly3ldmj6uxeea5n3fz3w2curfmw5647mwzlbvobd4j76g7cq
    ONTO initial
{
  CREATE TYPE default::Account {
      CREATE PROPERTY avatar_uri: std::str;
      CREATE PROPERTY created_at: std::datetime;
      CREATE REQUIRED PROPERTY email: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::regexp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
      };
      CREATE PROPERTY status: std::str {
          SET default := 'created';
          CREATE CONSTRAINT std::one_of('created', 'active', 'suspended', 'deleted');
      };
      CREATE PROPERTY status_changed: std::datetime;
      CREATE PROPERTY status_description: std::str;
      CREATE REQUIRED PROPERTY username: std::str;
  };
  CREATE TYPE default::Password {
      CREATE REQUIRED LINK account: default::Account;
      CREATE REQUIRED PROPERTY email: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::regexp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
      };
      CREATE PROPERTY failed_attempts: std::int16;
      CREATE PROPERTY last_failed_attempt: std::datetime;
      CREATE PROPERTY last_used: std::datetime;
      CREATE REQUIRED PROPERTY password: std::str {
          CREATE CONSTRAINT std::min_len_value(8);
      };
  };
};
