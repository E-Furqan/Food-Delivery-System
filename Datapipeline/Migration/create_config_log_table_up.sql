CREATE TABLE IF NOT EXISTS LogConfig (
    logId INT PRIMARY KEY,
    ConfigId INT,
    FOREIGN KEY (logId) REFERENCES Log(logId),
    FOREIGN KEY (ConfigId) REFERENCES Configuration(ConfigId)  
);