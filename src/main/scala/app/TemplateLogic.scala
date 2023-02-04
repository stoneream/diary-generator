package app

import org.joda.time.DateTime

object TemplateLogic {
  def render(template: String, now: DateTime): String = {
    template.replaceAll("%TODAY%", Format.ymd(now))
  }
}
