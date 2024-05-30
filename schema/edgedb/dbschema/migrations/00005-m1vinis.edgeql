CREATE MIGRATION m1vinisufbbduudedrobtmqlno27epjvrtbkjl6jy4jdvbin4pscda
    ONTO m1ztcgm6vgygcxx5x2jrhbvojrduqlkngqyntlakup7z3riuwl4rpa
{
  ALTER TYPE default::Password {
      ALTER PROPERTY failed_attempts {
          SET default := 0;
      };
  };
};
