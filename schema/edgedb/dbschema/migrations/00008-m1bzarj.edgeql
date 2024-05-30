CREATE MIGRATION m1bzarjufwic4s4jpbr7lhkzrhlkwvruside45hwwetri6i3neymda
    ONTO m1v5fr4o34ukw26znuyxv5fg5hlevr7nx4c534zazmcwkcqwue42ga
{
  ALTER TYPE default::Token {
      CREATE REQUIRED PROPERTY expires_in: std::datetime {
          SET REQUIRED USING (<std::datetime>{});
      };
  };
};
