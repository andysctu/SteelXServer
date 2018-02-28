insert into users values (6, 'JasonL', 'JasonL', 'JasonL', 1, 'Trainee', 0);
insert into owns (uid, eid) select 6, eid from owns where uid = 2;
insert into mechs (uid, arms, legs, core, head, weapon1l, weapon1r, weapon2l, weapon2r, booster, isprimary) select 6, arms, legs, core, head, weapon1l, weapon1r, weapon2l, weapon2r, booster, isprimary from mechs where uid = 2;
