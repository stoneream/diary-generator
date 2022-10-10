package app

import org.joda.time.DateTime
import org.joda.time.format.DateTimeFormat

object Format {
  def ymd(dt: DateTime): String = DateTimeFormat.forPattern("yyyy-MM-dd").print(dt)
}
