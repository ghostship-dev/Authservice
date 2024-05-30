CREATE MIGRATION m13zya64j63rrb5g5jykylxuku7ztn7cw32nhqiiatplx4fl552ffa
    ONTO m1tpz7fg4mgk6b4gzprk6tlhyx6e2rg6rdar6hfx3lx33z5wld5zeq
{
  ALTER TYPE default::Account {
      ALTER PROPERTY otp_state {
          SET default := 'disabled';
          DROP CONSTRAINT std::one_of('disabled', 'enabled', 'verifing');
      };
  };
};
