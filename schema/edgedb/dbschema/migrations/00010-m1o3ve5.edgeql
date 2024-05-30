CREATE MIGRATION m1o3ve57xneu4mjahekcwx6n36dnxivqoc4sq6ye3jzqhcybmilvyq
    ONTO m1jtqy3tpajpmh4dg4ujf6eawi3hcsfwyk6vpnqvzn7xgg5usz6xla
{
  ALTER TYPE default::Account {
      CREATE PROPERTY otp_secret: std::str;
      CREATE PROPERTY otp_state: std::str {
          CREATE CONSTRAINT std::one_of('disabled', 'enabled', 'verifing');
      };
  };
  CREATE TYPE default::OTP {
      CREATE REQUIRED PROPERTY secret: std::str;
      CREATE INDEX ON (.secret);
      CREATE REQUIRED LINK account: default::Account;
      CREATE REQUIRED PROPERTY backup_codes: array<std::str>;
  };
};
