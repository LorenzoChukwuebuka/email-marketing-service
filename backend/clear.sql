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


-- Version 1: Conservative approach (only drops custom indexes)
DO $$
DECLARE     
    table_record RECORD;     
    index_record RECORD;
    truncate_statement TEXT := ''; 
BEGIN     
    -- Disable all foreign key constraints     
    EXECUTE 'SET session_replication_role = replica;';          
    
    -- Drop all custom indexes (excluding primary keys, unique constraints, and system indexes)
    FOR index_record IN         
        SELECT schemaname, indexname, tablename
        FROM pg_indexes         
        WHERE schemaname = 'public'         
        AND indexname NOT LIKE '%_pkey'  -- Exclude primary key indexes
        AND indexname NOT LIKE '%_key'   -- Exclude unique constraint indexes
        AND indexname NOT IN (
            -- Exclude indexes that are part of constraints
            SELECT i.relname 
            FROM pg_constraint c
            JOIN pg_class i ON i.oid = c.conindid
            WHERE c.contype IN ('p', 'u', 'f')  -- primary, unique, foreign key
        )
    LOOP         
        BEGIN
            EXECUTE 'DROP INDEX IF EXISTS ' || quote_ident(index_record.schemaname) || '.' || quote_ident(index_record.indexname) || ';';             
            RAISE NOTICE 'Dropped index: %.%', index_record.schemaname, index_record.indexname;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Could not drop index %.%: %', index_record.schemaname, index_record.indexname, SQLERRM;
        END;
    END LOOP;          
    
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
    
    RAISE NOTICE 'Database reset complete: all custom indexes dropped and tables emptied'; 
END $$;

-- Version 2: Comprehensive cleanup (includes views, functions, sequences, etc.)
/*
DO $ 
DECLARE     
    record_item RECORD;     
    truncate_statement TEXT := ''; 
BEGIN     
    -- Disable all foreign key constraints     
    EXECUTE 'SET session_replication_role = replica;';          
    
    -- Drop all views
    FOR record_item IN 
        SELECT table_name 
        FROM information_schema.views 
        WHERE table_schema = 'public'
    LOOP
        EXECUTE 'DROP VIEW IF EXISTS ' || quote_ident(record_item.table_name) || ' CASCADE;';
        RAISE NOTICE 'Dropped view: %', record_item.table_name;
    END LOOP;
    
    -- Drop all custom functions
    FOR record_item IN 
        SELECT routine_name, routine_schema
        FROM information_schema.routines 
        WHERE routine_schema = 'public' 
        AND routine_type = 'FUNCTION'
    LOOP
        EXECUTE 'DROP FUNCTION IF EXISTS ' || quote_ident(record_item.routine_schema) || '.' || quote_ident(record_item.routine_name) || ' CASCADE;';
        RAISE NOTICE 'Dropped function: %', record_item.routine_name;
    END LOOP;
    
    -- Drop all custom indexes
    FOR record_item IN         
        SELECT schemaname, indexname
        FROM pg_indexes         
        WHERE schemaname = 'public'         
        AND indexname NOT LIKE '%_pkey'
        AND indexname NOT LIKE '%_key'
        AND indexname NOT IN (
            SELECT i.relname 
            FROM pg_constraint c
            JOIN pg_class i ON i.oid = c.conindid
            WHERE c.contype IN ('p', 'u', 'f')
        )
    LOOP         
        EXECUTE 'DROP INDEX IF EXISTS ' || quote_ident(record_item.schemaname) || '.' || quote_ident(record_item.indexname) || ';';             
        RAISE NOTICE 'Dropped index: %', record_item.indexname;
    END LOOP;          
    
    -- Build TRUNCATE statement for all tables     
    FOR record_item IN         
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
        truncate_statement := truncate_statement || quote_ident(record_item.table_name);     
    END LOOP;          
    
    -- Execute the truncate if we have tables     
    IF truncate_statement != '' THEN         
        EXECUTE 'TRUNCATE TABLE ' || truncate_statement || ' RESTART IDENTITY CASCADE;';         
        RAISE NOTICE 'Truncated tables: %', truncate_statement;     
    END IF;          
    
    -- Drop orphaned sequences
    FOR record_item IN
        SELECT sequence_name 
        FROM information_schema.sequences 
        WHERE sequence_schema = 'public'
        AND sequence_name NOT IN (
            SELECT pg_get_serial_sequence(schemaname||'.'||tablename, column_name) 
            FROM information_schema.columns c
            JOIN information_schema.tables t ON c.table_name = t.table_name
            WHERE pg_get_serial_sequence(schemaname||'.'||tablename, column_name) IS NOT NULL
        )
    LOOP
        EXECUTE 'DROP SEQUENCE IF EXISTS ' || quote_ident(record_item.sequence_name) || ';';
        RAISE NOTICE 'Dropped orphaned sequence: %', record_item.sequence_name;
    END LOOP;
    
    -- Re-enable foreign key constraints     
    EXECUTE 'SET session_replication_role = DEFAULT;';          
    
    RAISE NOTICE 'Complete database cleanup finished: views, functions, indexes dropped and tables emptied'; 
END $;
*/