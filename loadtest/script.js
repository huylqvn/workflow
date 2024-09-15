
import http, { file } from 'k6/http';
import { check, sleep } from 'k6';

const apiEndpoint = 'http://127.0.0.1:8088';

export const options = {
    iterations: 10,
    vus: 10,
    duration: '1m',
};


export default function () {
    const resp = http.post(apiEndpoint + '/v1/workflow',
        `
        {
            "name":"abc",
            "description":"test",
            "tasks":[
                {
                    "name":"task1",
                    "message":"xin",
                    "job_type":"period",
                    "job_time_value":"1s"
                },
                {
                    "name":"task2",
                    "message":"chao",
                    "job_type":"period",
                    "job_time_value":"1s"
                },
                {
                    "name":"task3",
                    "message":"viet",
                    "job_type":"period",
                    "job_time_value":"1s"
                },
                {
                    "name":"task4",
                    "message":"nam",
                    "job_type":"period",
                    "job_time_value":"1s"
                }
            ]
        }
        `
    )

    check(resp, {
        'status is 200': (r) => r.status === 201
    })
    sleep(1);
}