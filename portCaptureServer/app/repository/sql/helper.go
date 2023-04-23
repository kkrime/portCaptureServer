package sql

import (
	"context"
	"database/sql"
	"portCaptureServer/app/entity"
)

type db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// common code
func savePort(ctx context.Context, db db, port *entity.Port) error {
	// 1. mark deleted any existing port with the same code,
	// this is for auditing reasons
	statement := `
        UPDATE
            ports
        SET
            deleted_at = now()
        WHERE
            primary_unloc = $1
        ;`

	_, err := db.ExecContext(ctx, statement, port.PrimaryUnloc)
	if err != nil {
		return err
	}

	// 2. add the port
	statement = `
        INSERT INTO
            ports
        (
            primary_unloc,
            name,
            code,
            city,
            country,
            coordinates,
            province,
            timezone
        )
        VALUES
        (
            $1,
            $2,
            $3,
            $4,
            $5,
            point($6,$7),
            $8,
            $9
        )
        ;`

	_, err = db.ExecContext(ctx, statement,
		port.PrimaryUnloc,
		port.Name,
		port.Code,
		port.City,
		port.Country,
		port.Coordinantes[0],
		port.Coordinantes[1],
		port.Province,
		port.Timezone,
	)
	if err != nil {
		return err
	}

	// 3. add alias
	for _, alias := range *port.Alias {
		statement = `
            INSERT INTO
                alias
            (
                port_id,
                name
            )
            VALUES
            (
                (
                    SELECT 
                        id 
                    FROM 
                        ports
                    WHERE 
                        primary_unloc = $1 AND 
                        deleted_at IS NULL
                ),
                $2
            )`

		_, err = db.ExecContext(ctx, statement, port.PrimaryUnloc, alias.Name)
		if err != nil {
			return err
		}

	}

	// 4. add regions
	for _, region := range *port.Regions {
		statement = `
            INSERT INTO
                regions
            (
                port_id,
                name
            )
            VALUES
            (
                (
                    SELECT 
                        id 
                    FROM 
                        ports
                    WHERE 
                        primary_unloc = $1 AND 
                        deleted_at IS NULL
                ),
                $2
            )`

		_, err = db.ExecContext(ctx, statement, port.PrimaryUnloc, region.Name)
		if err != nil {
			return err
		}
	}

	// 5. add unlocs
	for _, unloc := range *port.Unlocs {
		statement = `
            INSERT INTO
                unlocs
            (
                port_id,
                name
            )
            VALUES
            (
                (
                    SELECT 
                        id 
                    FROM 
                        ports
                    WHERE 
                        primary_unloc = $1 AND 
                        deleted_at IS NULL
                ),
                $2
            )`

		_, err = db.ExecContext(ctx, statement, port.PrimaryUnloc, unloc.Name)
		if err != nil {
			return err
		}
	}

	return nil

}
