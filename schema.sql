-- users テーブルの作成
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- characters テーブルの作成
CREATE TABLE IF NOT EXISTS characters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    rarity INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- gacha_probabilities テーブルの作成
CREATE TABLE IF NOT EXISTS gacha_probabilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    character_id INT NOT NULL,
    probability FLOAT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id)
) ENGINE=InnoDB;

-- user_characters テーブルの作成
CREATE TABLE IF NOT EXISTS user_characters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    character_id INT NOT NULL,
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (character_id) REFERENCES characters(id)
) ENGINE=InnoDB;

-- キャラクターの初期データ
INSERT INTO characters (name, rarity) VALUES
('Warrior', 1),
('Mage', 2),
('Archer', 3),
('Knight', 4),
('Dragon', 5);

-- ガチャ確率の初期データ
INSERT INTO gacha_probabilities (character_id, probability) VALUES
(1, 0.4),  -- Warrior
(2, 0.3),  -- Mage
(3, 0.2),  -- Archer
(4, 0.08), -- Knight
(5, 0.02); -- Dragon
