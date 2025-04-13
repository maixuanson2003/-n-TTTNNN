INSERT IGNORE  INTO songs (id,
  name_song, description, release_day, create_day, update_day,
  point, like_amount, status, country_id, listen_amout, album_id, song_resource
) VALUES
  (1,'Cham em Một Đời', 'Bài hát về cha mẹ', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/chamemmotdoi.mp3'),
  (2,'Chín Tầng Mây', 'Bài hát nhẹ nhàng', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/chintangmay.mp3'),
  (3,'Lỡ Yêu', 'Ca khúc ballad tình cảm', NOW(), NOW(), NOW(), 0, 0, "release", 2, 0, null, 'http://localhost:8080/music/loyeu.mp3');
