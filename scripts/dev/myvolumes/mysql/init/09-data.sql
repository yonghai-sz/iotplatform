USE `myexampledb`;

SET names utf8mb4;






BEGIN;

INSERT INTO `role` ( `role_type`, `role_name`, `created_at`) VALUES ( 'SuperAdmin', 'super administrator', NOW(3) );


INSERT INTO `user` ( `username`, `password`, `salt`, `role_id`, `created_at`)
VALUES	( 'super-admin', 'b7d6c2ee39faf9179ef1ae83d2dbb6842da3915a8b3370f376c2675c953fc67fa27e2d411f620be21fc1f80b338aeb7bd1ea58162b7ceb6c1aac991b6a5ef65a',
    'abc456',  1, NOW(3) );

INSERT INTO `menu` (`id`, `parent_id`, `title`, `has_child`, `menu_key`) VALUES
  (54, 0, '视频中心', 'N', 'videoCenter'),
  (59, 0, '报警记录', 'N', 'alarm');

INSERT INTO `role_menu` (`role_id`, `menu_id`) VALUES (1, 54), (1, 59);



COMMIT;