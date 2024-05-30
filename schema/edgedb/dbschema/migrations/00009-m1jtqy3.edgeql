CREATE MIGRATION m1jtqy3tpajpmh4dg4ujf6eawi3hcsfwyk6vpnqvzn7xgg5usz6xla
    ONTO m1bzarjufwic4s4jpbr7lhkzrhlkwvruside45hwwetri6i3neymda
{
  ALTER TYPE default::Token {
      ALTER PROPERTY expires_in {
          RENAME TO expires_at;
      };
  };
};
