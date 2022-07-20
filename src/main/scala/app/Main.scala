package app

import scopt.OParser

object Main extends App {
  val parser = {
    val builder = OParser.builder[Config]
    import builder._
    OParser.sequence(
      opt[String]("target-path").required().action((targetPath, config) => config.copy(targetPath = targetPath)),
      cmd("init")
        .action((_, config) => config.copy(mode = "init"))
        .children(
          opt[String]("template-path").required().action((templatePath, config) => config.copy(templatePathOpt = Some(templatePath)))
        ),
      cmd("archive")
        .action((_, config) => config.copy(mode = "archive"))
        .children(
          opt[String]("starts-with").required().action((startsWith, config) => config.copy(startsWithOpt = Some(startsWith)))
        )
    )
  }

  OParser
    .parse(parser, args, Config())
    .fold {
      // do nothing
    } { config =>
      (config.mode, config.templatePathOpt, config.startsWithOpt) match {
        case ("init", Some(_), _) => ???
        case ("archive", _, Some(_)) => ???
        case _ => ???
      }
    }
}
