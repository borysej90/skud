CREATE VIEW employee_current_nodes AS
SELECT employee_id, node_id
FROM (
         SELECT employee_id, node_id, MAX(t1.access_time) AS a_time
         FROM access_log AS T1
                  JOIN access_log AS t2 USING (employee_id, node_id)
         WHERE t1.access + t2.access = 2
           AND t1.exited - t2.exited = -1
         GROUP BY employee_id, node_id
         HAVING a_time > MAX(t2.access_time)
         UNION
         SELECT employee_id, node_id, MAX(access_time) AS access_time
         FROM access_log
         WHERE access = 1
         GROUP BY employee_id, node_id
         HAVING COUNT(id) = 1
     ) AS sub
ORDER BY employee_id, a_time DESC;
