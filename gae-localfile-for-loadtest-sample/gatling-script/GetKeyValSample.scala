/*
 * Copyright 2011-2019 GatlingCorp (https://gatling.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package computerdatabase.advanced

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

// sh bin/gatling.sh
class GetKeyValSample extends Simulation {

  val httpProtocol = http
    .baseUrl("http://localhost:8080")
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  // Now, we can write the scenario as a composition
  val users = scenario("GetKeyValSample")
    .exec(http("request")
      .get("/values/abc") // OK
      //.get("/values/xxx") // NG
      .check(bodyString.saveAs("bodyString")))
    .exec( session => {
      println( "Result bodyString: " )
      println( session( "bodyString" ).as[String] )
      session
    })

  setUp(
    users.inject(
      rampUsersPerSec(1) to (3) during (3 seconds),
      constantUsersPerSec(3) during(5 seconds)
      //rampUsersPerSec(1) to (100) during (100 seconds),
      //constantUsersPerSec(100) during(10 seconds)
    ).protocols(httpProtocol)
  )
}
