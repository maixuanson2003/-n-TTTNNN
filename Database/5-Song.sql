INSERT IGNORE  INTO songs (id,
  name_song, description, release_day, create_day, update_day,
  point, like_amount, status, country_id, listen_amout, album_id, song_resource
) VALUES
  (1,'Cham em Một Đời', 'Bài hát về cha mẹ', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/chamemmotdoi.mp3'),
  (2,'Chín Tầng Mây', 'Bài hát nhẹ nhàng', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/chintangmay.mp3'),
  (3,'Lỡ Yêu', 'Ca khúc ballad tình cảm', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/loyeu.mp3'),
  (4, 'Da Hool - Meet her at the love parade (RBX EDIT)', 'Bản remix sôi động', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/Da Hool - Meet her at the love parade (RBX EDIT).wav'),
  (5, 'Don''t Let Me Down (feat. Daya) (Hipst3r Edit)', 'Bản EDM remix nổi bật', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/Don''t Let Me Down (feat. Daya) (Hipst3r Edit).wav'),
  (6,'Ely Oaks - Running Around (Deny DnB Remix) Master', 'Bản DnB remix mạnh mẽ', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/Ely Oaks - Running Around (Deny DnB Remix) Master.wav'),
  (7,'Greedy Remix V3 Extended', 'Phiên bản remix Greedy mở rộng', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/Greedy Remix V3 Extended.wav'),
  (8,'POPULAR - MINDLOCO REMIX (EDIT)', 'Remix hiện đại và bắt tai', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/POPULAR - MINDLOCO REMIX (EDIT).wav'),
  (9,'Purpura Violaceum x Habits Tove Lo', 'Kết hợp độc đáo giữa Purpura và Tove Lo', NOW(), NOW(), NOW(), 0, 0, 'release', 2, 0, NULL, 'http://localhost:8080/music/Purpura Violaceum x Habits Tove Lo.mp3');
