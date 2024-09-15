Preparation
Please do coding and prepare design docs before the pair session.
We love Go and Postgres. But sure you are free to use any language and database.
Pair session
We will first discuss your solution.
A code review session comes after that.
And then let’s talk about further issues and improvements.
Problem brief
Our Workflow product has a feature called the Time delay. It allows users to specify exactly how long the system should wait before moving to the next step, i.e. run the scheduled actions, in their workflow. We support specifying timing in the following ways:
A certain period of time. It means waiting for X minutes, hours, or days before doing the job.
A certain day of the week. For example, if the user selects Monday and Wednesday, but the current day is Tuesday, the system will delay the job until Wednesday because it is the first available day of the week match.
A specific day of year. The system will wait until the date matches the user’s inputted date.

Depending on the traffic, there can be millions of new/due timers every hour. The scheduled time can be varied, from a few minutes or hours to months or even years into the future, and they must be delivered reliably to ensure a good experience for the users.
We need to build a job scheduling system to support the above feature with the following challenges:
Reliability. Timers must be recovered if the system fails or restarts. It should guarantee at-least-once delivery of every single job.
Scalability. The scheduling system should be horizontally scalable, to be able to handle ~5M new/due timers per hour at peak. 
Timing accuracy. In the workflow system, it is required that jobs run within seconds of their scheduled time. The scheduler should have a p95 scheduling deviation below 10 seconds.
Flexibility. The user is able to change value or cancel timers on the fly. 

Tasks
System design
Please provides proper design documents including:
A high level diagram showing how components wiring together, and database schema.
A document describes your solution and the technologies.
The design should be able to express the capabilities of tackling following challenges:
The system’s capacity should cover 5M new timers and 5M due timers per hour. The traffic can be doubled anywhen in near future.
Cost should be considered of course. Some component such as worker should be autoscalable.
Last but not least: fairness.
As we maintain the same price for all users, we won’t expect a few “whales” eat up all the capacity. The system should also give priority to other small business owners with small audience list.
To maximize our paid conversion rate, we should give the best experience to new signed-up users who’s testing the workflow feature with a short time delay.
Implementation
The implementation should be able to provide:
A proof that the concept is feasible.
A way to test and benchmark.
For further discussion during pair session
We need proper monitoring for the performance and latency of the system (also refer to the Timing Accuracy challenge above). So let’s discuss to define the monitoring metrics.
Notes for the scope
Your system will handle the due timers and output the jobs into a queue, from which the actual executor will pull the jobs.
The scope of this assignment only requires to design and implement the schedule part, you don't have to care about the actual job execution.
To show us your skill, we require you to design the system without using or extending any Cloud or Open-source Job scheduler.
The implementation doesn’t need to be production ready but to show us how you architect a clean code base that well matches with the system design and feature requirements.

If you have any questions please don’t hesitate to reach out!
Attention: please ensure that your submitted document and code have no references to our branding.