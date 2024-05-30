CREATE MIGRATION m1ebs7a66ao6qmaeo2yo6p3w6bwnujdn3hliyy2eehlp4f6sh24drq
    ONTO m1l2ryn5vglpivib2wpdix6ofogwjck65unb6t5varn4fiwdw44boq
{
  ALTER TYPE default::Account {
      ALTER PROPERTY otp_state {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
