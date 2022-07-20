package app

case class Config(
    mode: String = "",
    targetPath: String = "",
    templatePathOpt: Option[String] = None,
    startsWithOpt: Option[String] = None
)
