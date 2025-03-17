<?php
$servername = "SRV_NAME";
$username = "USERNAME";
$password = "PASSWORD";
$dbname = "DB_NAME";

try {
    // Connect to MySQL using PDO
    $pdo = new PDO("mysql:host=$servername;dbname=$dbname;charset=utf8mb4", $username, $password);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    // Get all tables
    $tables = [];
    $stmt = $pdo->query("SHOW TABLES");
    while ($row = $stmt->fetch(PDO::FETCH_NUM)) {
        $tables[] = $row[0];
    }

    // Dump each table
    foreach ($tables as $table) {
        echo "\n--- Dumping table: $table ---\n";
        
        // Get table structure
        $stmt = $pdo->query("SHOW CREATE TABLE `$table`");
        $createTable = $stmt->fetch(PDO::FETCH_ASSOC);
        echo $createTable['Create Table'] . ";\n";

        // Get table data
        $stmt = $pdo->query("SELECT * FROM `$table`");
        while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $escapedValues = array_map(fn($val) => $pdo->quote($val), $row);
            echo "INSERT INTO `$table` VALUES (" . implode(",", $escapedValues) . ");\n";
        }
    }

} catch (PDOException $e) {
    die("Connection failed: " . $e->getMessage());
}
// check for available php module with : phm -m and use the sql driver accordingly
?>
