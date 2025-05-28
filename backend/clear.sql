DO $$
DECLARE
    table_record RECORD;
    truncate_statement TEXT := '';
BEGIN
    -- Disable all foreign key constraints
    EXECUTE 'SET session_replication_role = replica;';
    
    -- Build TRUNCATE statement for all tables
    FOR table_record IN
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        AND table_name NOT LIKE 'pg_%'
        AND table_name NOT LIKE 'sql_%'
    LOOP
        IF truncate_statement != '' THEN
            truncate_statement := truncate_statement || ', ';
        END IF;
        truncate_statement := truncate_statement || quote_ident(table_record.table_name);
    END LOOP;
    
    -- Execute the truncate if we have tables
    IF truncate_statement != '' THEN
        EXECUTE 'TRUNCATE TABLE ' || truncate_statement || ' RESTART IDENTITY CASCADE;';
        RAISE NOTICE 'Truncated tables: %', truncate_statement;
    ELSE
        RAISE NOTICE 'No tables found to truncate';
    END IF;
    
    -- Re-enable foreign key constraints
    EXECUTE 'SET session_replication_role = DEFAULT;';
    
    RAISE NOTICE 'All tables have been emptied successfully';
END $$;