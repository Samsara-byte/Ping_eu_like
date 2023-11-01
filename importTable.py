import csv
import sqlite3

# Connect to the SQLite database
conn = sqlite3.connect('ip-location.db')
cursor = conn.cursor()

# Create the table
create_table_query = '''
CREATE TABLE `proxy`(
	`ip_from` INT,
	`ip_to` INT,
	`proxy_type` VARCHAR(3),
	`country_code` TEXT,
	`country_name` VARCHAR(64),
	`region_name` VARCHAR(128),
	`city_name` VARCHAR(128),
	`isp` VARCHAR(256),
	`domain` VARCHAR(128),
	`usage_type` VARCHAR(11),
	`asn` INT,
	`as_name` VARCHAR(256)
) 
'''
cursor.execute(create_table_query)

# Read the CSV file and insert data into the table
csv_filename = 'IP2PROXY-LITE-PX7.CSV'

with open(csv_filename, 'r') as csvfile:
    csvreader = csv.reader(csvfile)
    next(csvreader) 
    
    for row in csvreader:
        insert_query = '''
        INSERT INTO proxy (ip_from, ip_to, proxy_type,country_code,country_name,region_name,city_name,isp,domain,usage_type, asn, as_name)
        VALUES (?, ?, ?, ?, ?,?, ?, ?, ?, ?,?,?)
        '''
        cursor.execute(insert_query, row)

# Commit the changes and close the database connection
conn.commit()
conn.close()
